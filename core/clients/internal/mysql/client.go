package mysql

import (
	"context"
	"strings"
	"sync"

	"github.com/mel0dys0ng/song/core/erlogs"
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

	elgSys := erlogs.New(erlogs.TypeSystem(), erlogs.Log(true), erlogs.Msgf("[mysql] %s"))
	config, err := newConfig(ctx, key, elgSys, opts)
	if err != nil {
		err.PanicL(ctx)
		return
	}

	logger, err := newLogger(ctx, elgSys, config)
	if err != nil {
		err.PanicL(ctx)
		return
	}

	client := &Client{key: key, config: config}
	gormConfig := &gorm.Config{
		Logger: logger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   client.config.TablePrefix,
			SingularTable: client.config.SingularTable,
		},
	}

	dialector := mysql.Open(client.config.Master)
	db, er := gorm.Open(dialector, gormConfig)
	if er != nil {
		elgSys.PanicL(ctx,
			erlogs.Msgv("gorm open failed"),
			erlogs.Content(er.Error()),
		)
		return
	}

	er = db.Use(newResolver(client.config))
	if er != nil {
		elgSys.PanicL(ctx,
			erlogs.Msgv("use plugin failed"),
			erlogs.Content(er.Error()),
		)
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
