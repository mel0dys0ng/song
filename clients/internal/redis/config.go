package redis

import (
	"context"
	"time"

	"github.com/mel0dys0ng/song/erlogs"
	"github.com/mel0dys0ng/song/vipers"
)

// Config 数据库配置
type Config struct {
	// debug模式，默认false
	Debug bool `mapstructure:"debug" json:"debug"`

	// IdentitySuffix
	IdentitySuffix string `mapstructure:"identitySuffix" json:"identitySuffix"`

	// MasterName The sentinel failover master name. sentinel failover option
	MasterName string `mapstructure:"masterName" json:"masterName"`

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string `mapstructure:"clientName" json:"clientName"`

	// 网络类型，tcp or unix，默认tcp
	Network string `mapstructure:"network" json:"network"`

	// 连接地址和端口, host:port addresses. 非cluster redis，取第一个addr
	Addrs []string `mapstructure:"addrs" json:"addrs"`

	// 连接账号
	Username string `mapstructure:"username" json:"username"`

	// 连接密码
	Password string `mapstructure:"password" json:"password"`

	// Sentinel连接账号
	SentinelUsername string `mapstructure:"sentinelUsername" json:"sentinelUsername"`

	// Sentinel连接密码
	SentinelPassword string `mapstructure:"sentinelPassword" json:"sentinelPassword"`

	// Database to be selected after connecting to the server.
	Database int `mapstructure:"database" json:"database"`

	// 命令执行失败时，最多重试多少次，默认为0即不重试
	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int `mapstructure:"maxRetries" json:"maxRetries"`

	// 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration `mapstructure:"minRetryBackoff" json:"minRetryBackoff"`

	// 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration `mapstructure:"maxRetryBackoff" json:"maxRetryBackoff"`

	// 连接建立超时时间，单位:毫秒
	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration `mapstructure:"dialTimeout" json:"dialTimeout"`

	// 读超时，单位:毫秒
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout time.Duration `mapstructure:"readTimeout" json:"readTimeout"`

	// 写超时，单位:毫秒
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout"`

	// 当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。单位：毫秒
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration `mapstructure:"poolTimeout" json:"poolTimeout"`

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that FIFO has slightly higher overhead compared to LIFO,
	// but it utils closing idle connections faster reducing the pool size.
	PoolFIFO bool `mapstructure:"poolFIFO" json:"poolFIFO"`

	// 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
	// Maximum number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int `mapstructure:"poolSize" json:"poolSize"`

	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConns int `mapstructure:"minIdleConns" json:"minIdleConns"`

	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConns int `mapstructure:"maxIdleConns" json:"maxIdleConns"`

	// Maximum number of active connections.
	MaxActiveConns int `mapstructure:"maxActiveConns" json:"maxActiveConns"`

	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	//
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	//
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime" json:"connMaxIdleTime"`

	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	//
	// Default is to not close idle connections.
	ConnMaxLifetime time.Duration `mapstructure:"connMinIdleTime" json:"connMinIdleTime"`

	// 空闲超时，单位:毫秒
	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration `mapstructure:"idleTimeout" json:"idleTimeout"`

	// The maximum number of retries before giving up. Command is retried
	// on network errors and MOVED/ASK redirects.
	// Default is 3 retries.
	MaxRedirects int `mapstructure:"maxRedirects" json:"maxRedirects"`

	// ContextTimeoutEnabled controls whether the client respects context timeouts and deadlines.
	// See https://redis.uptrace.dev/guide/go-redis-debugging.html#timeouts
	ContextTimeoutEnabled bool `mapstructure:"contextTimeoutEnabled" json:"contextTimeoutEnabled"`

	// Enables read-only commands on slave nodes.
	ReadOnly bool `mapstructure:"readOnly" json:"readOnly"`

	// Allows routing read-only commands to the closest master or slave node.
	// It automatically enables ReadOnly.
	RouteByLatency bool `mapstructure:"routeByLatency" json:"routeByLatency"`

	// Allows routing read-only commands to the random master or slave node.
	// It automatically enables ReadOnly.
	RouteRandomly bool `mapstructure:"routeRandomly" json:"routeRandomly"`
}

func NewConfig(ctx context.Context, key string, elg erlogs.ErLogInterface, opts []Option) (
	config *Config, err erlogs.ErLogInterface) {

	config = &Config{
		Debug:                 DefauktDebug,
		MasterName:            DefaultMasterName,
		ClientName:            DefaultClientName,
		Network:               DefaultNetwork,
		Username:              DefaultUsername,
		Password:              DefaultPassword,
		Database:              DefaultDatabase,
		MaxRetries:            DefaultMaxRetries,
		MinRetryBackoff:       DefaultMinRetryBackoff,
		MaxRetryBackoff:       DefaultMaxRetryBackoff,
		DialTimeout:           DefaultDialTimeout,
		ReadTimeout:           DefaultReadTimeout,
		WriteTimeout:          DefaultWriteTimeout,
		PoolFIFO:              DefaultPoolFIFO,
		PoolSize:              DefaultPoolSize,
		PoolTimeout:           DefaultPoolTimeout,
		MinIdleConns:          DefaultMinIdleConns,
		MaxIdleConns:          DefaultMaxIdleConns,
		ConnMaxIdleTime:       DefaultConnMaxIdleTime,
		ConnMaxLifetime:       DefaultConnMaxLifetime,
		IdleTimeout:           DefaultIdleTimeout,
		MaxRedirects:          DefaultMaxRedirects,
		ContextTimeoutEnabled: DefaultContextTimeoutEnabled,
		ReadOnly:              DefaultReadOnly,
		RouteRandomly:         DefaultRouteRandomly,
		RouteByLatency:        DefaultRouteByLatency,
	}

	options, err := config.LoadOptions(ctx, key, elg)
	if err != nil {
		return
	}

	opts = append(options, opts...)
	for _, v := range opts {
		v.Func(config)
	}

	err = config.check(ctx, elg)

	return
}

// LoadOptions 加载配置并返回Options
func (c *Config) LoadOptions(ctx context.Context, key string, elg erlogs.ErLogInterface) (
	opts []Option, err erlogs.ErLogInterface) {

	values := &Config{}
	er := vipers.UnmarshalKey(key, values, nil)
	if er != nil {
		err = elg.PanicE(ctx,
			erlogs.Msgv("failed to load options: UnmarshalKey error"),
			erlogs.Content(er.Error()),
		)
		return
	}

	if values.IdentitySuffix != DefaultIdentifySufix {
		opts = append(opts, IdentitySuffix(values.IdentitySuffix))
	}

	if values.MasterName != DefaultMasterName {
		opts = append(opts, MasterName(values.MasterName))
	}

	if values.ClientName != DefaultClientName {
		opts = append(opts, ClientName(values.ClientName))
	}

	if values.Network != DefaultNetwork {
		opts = append(opts, Network(values.Network))
	}

	if len(values.Addrs) != 1 || (len(values.Addrs) == 1 && values.Addrs[0] != DefaultAddr) {
		opts = append(opts, Addrs(values.Addrs))
	}

	if values.SentinelUsername != DefaultSentinelUsername {
		opts = append(opts, SentinelUsername(values.SentinelUsername))
	}

	if values.SentinelPassword != DefaultSentinelPassword {
		opts = append(opts, SentinelPassword(values.SentinelPassword))
	}

	if values.Username != DefaultUsername {
		opts = append(opts, Username(values.Username))
	}

	if values.Password != DefaultPassword {
		opts = append(opts, Password(values.Password))
	}

	if values.Database != DefaultDatabase {
		opts = append(opts, Database(values.Database))
	}

	if values.MaxRetries != DefaultMaxRetries {
		opts = append(opts, MaxRetries(values.MaxRetries))
	}

	if values.MinRetryBackoff != DefaultMinRetryBackoff {
		opts = append(opts, MinRetryBackoff(values.MinRetryBackoff))
	}

	if values.MaxRetryBackoff != DefaultMaxRetries {
		opts = append(opts, MaxRetryBackoff(values.MaxRetryBackoff))
	}

	if values.MaxRedirects != DefaultMaxRedirects {
		opts = append(opts, MaxRedirects(values.MaxRedirects))
	}

	if values.MaxIdleConns != DefaultMaxIdleConns {
		opts = append(opts, MaxIdleConns(values.MaxIdleConns))
	}

	if values.MinIdleConns != DefaultMinIdleConns {
		opts = append(opts, MinIdleConns(values.MinIdleConns))
	}

	if values.PoolSize != DefaultPoolSize {
		opts = append(opts, PoolSize(values.PoolSize))
	}

	if values.ConnMaxIdleTime != DefaultConnMaxIdleTime {
		opts = append(opts, ConnMaxIdleTime(values.ConnMaxIdleTime))
	}

	if values.ConnMaxLifetime != DefaultConnMaxLifetime {
		opts = append(opts, ConnMaxLifetime(values.ConnMaxLifetime))
	}

	if values.IdleTimeout != DefaultIdleTimeout {
		opts = append(opts, IdleTimeout(values.IdleTimeout))
	}

	if values.ReadTimeout != DefaultReadTimeout {
		opts = append(opts, ReadTimeout(values.ReadTimeout))
	}

	if values.WriteTimeout != DefaultWriteTimeout {
		opts = append(opts, WriteTimeout(values.WriteTimeout))
	}

	if values.DialTimeout != DefaultDialTimeout {
		opts = append(opts, DialTimeout(values.DialTimeout))
	}

	if values.PoolTimeout != DefaultPoolTimeout {
		opts = append(opts, PoolTimeout(values.PoolTimeout))
	}

	return
}

func (c *Config) check(ctx context.Context, elg erlogs.ErLogInterface) (err erlogs.ErLogInterface) {
	if len(c.Addrs) == 0 {
		err = elg.PanicE(ctx, erlogs.Msgv("config addrs empty"))
		return
	}

	return
}
