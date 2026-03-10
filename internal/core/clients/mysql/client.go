package mysql

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

type (
	Client struct {
		config *Config
		key    string
		db     *gorm.DB
	}
)

var (
	clients = &sync.Map{}
)

// CreateClient 返回MySQL数据库链接，支持读写分离、负载均衡、连接池
// @key string database config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @opts []Option 自定义配置选项
func CreateClient(ctx context.Context, name, key string, opts ...Option) (res *Client) {
	mk := key
	if len(name) > 0 {
		mk = strings.Join([]string{key, name}, "-")
	}

	if v, ok := clients.Load(mk); ok {
		if client, ok := v.(*Client); ok {
			return client
		}
	}

	config, err := newConfig(ctx, key, opts)
	if err != nil {
		erlogs.Convert(err).Wrap("failed to create client: newConfig error").Options(BaseELOptions()).RecordLog(ctx)
		return
	}

	logger, err := newLogger(config)
	if err != nil {
		erlogs.Convert(err).Wrap("failed to create client: newLogger error").Options(BaseELOptions()).PanicLog(ctx)
		return
	}

	client := &Client{key: key, config: config}
	gormConfig := &gorm.Config{
		Logger: logger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   client.config.TablePrefix,
			SingularTable: *client.config.SingularTable,
		},
	}

	dialector := mysql.Open(client.config.Master)
	db, er := gorm.Open(dialector, gormConfig)
	if er != nil {
		erlogs.Convert(er).Wrap("failed to create client: gorm open error").Options(BaseELOptions()).PanicLog(ctx)
		return
	}

	er = db.Use(newResolver(client.config))
	if er != nil {
		erlogs.Convert(er).Wrap("failed to create client: use plugin error").Options(BaseELOptions()).PanicLog(ctx)
		return
	}

	client.db = db
	clients.Store(mk, client)

	return client
}

// Key return the database config key
func (c *Client) Key() string {
	return c.key
}

// Master database master, for write connection
func (c *Client) Master() *gorm.DB {
	return c.db.Clauses(dbresolver.Write)
}

// Slave database slave, for read connection
func (c *Client) Slave() *gorm.DB {
	return c.db.Clauses(dbresolver.Read)
}

// BeginTransaction 开始一个新的事务
func (c *Client) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return c.Master().Begin(opts...)
}

// DB 根据参数决定返回事务数据库实例还是主从数据库实例
// 当传入的事务实例不为空时，直接返回该事务实例
// 当useSlave为true时，返回从库实例用于读操作；否则返回主库实例用于写操作
// 参数:
//
//	tx: GORM事务数据库实例，如果非空则直接返回此实例
//	useSlave: 布尔值，指定是否使用从库进行查询操作
//
// 返回值:
//
//	*gorm.DB: 根据条件选择的GORM数据库实例
func (c *Client) DB(tx *gorm.DB, useSlave bool) *gorm.DB {
	if tx != nil {
		return tx
	}

	if useSlave {
		return c.Slave()
	}

	return c.Master()
}

// Close closes the database connection
func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping tests the database connection
func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
