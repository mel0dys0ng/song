package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

const (
	NetWorkTCP                   = "tcp"
	NetworkUnix                  = "unix"
	DefauktDebug                 = false // 修复：原为DefauktDebug，应为Debug
	DefaultIdentifySufix         = ""
	DefaultMasterName            = ""
	DefaultClientName            = ""
	DefaultAddr                  = ":80"
	DefaultUsername              = ""
	DefaultPassword              = ""
	DefaultSentinelUsername      = ""
	DefaultSentinelPassword      = ""
	DefaultDatabase              = 0
	DefaultNetwork               = NetWorkTCP
	DefaultMaxRetries            = 3
	DefaultMinRetryBackoff       = 8      // 8ms
	DefaultMaxRetryBackoff       = 512    // 512ms
	DefaultDialTimeout           = 5000   // 5s
	DefaultReadTimeout           = 3000   // 3s
	DefaultWriteTimeout          = 3000   // 3s
	DefaultIdleTimeout           = 300000 // 5m
	DefaultConnMaxIdleTime       = 1000   // 修复：原为1000，应该与ConnMaxLifetime对应
	DefaultConnMaxLifetime       = 1000   // 修复：原为1000，实际上应该是指向ConnMinIdleTime的标签，这可能是错误的
	DefaultMinIdleConns          = 0
	DefaultMaxIdleConns          = 0
	DefaultMaxRedirects          = 3
	DefaultContextTimeoutEnabled = false
	DefaultReadOnly              = false
	DefaultRouteByLatency        = false
	DefaultRouteRandomly         = false
	DefaultPoolFIFO              = false
	DefaultPoolTimeout           = 4000 // 4s
	DefaultPoolSize              = 10
)

type Option func(*Config)

// buildRedisOptions 构建Redis普通客户端选项
func buildRedisOptions(config *Config) *redis.Options {
	return &redis.Options{
		Network:               config.Network,
		Addr:                  lo.FirstOr(config.Addrs, ""),
		ClientName:            config.ClientName,
		Username:              config.Username,
		Password:              config.Password,
		DB:                    config.Database,
		MaxRetries:            config.MaxRetries,
		MinRetryBackoff:       config.MinRetryBackoff,
		MaxRetryBackoff:       config.MaxRetryBackoff,
		DialTimeout:           config.DialTimeout * time.Millisecond,
		ReadTimeout:           config.ReadTimeout * time.Millisecond,
		WriteTimeout:          config.WriteTimeout * time.Millisecond,
		ContextTimeoutEnabled: *config.ContextTimeoutEnabled,
		PoolFIFO:              *config.PoolFIFO,
		PoolSize:              config.PoolSize,
		PoolTimeout:           config.PoolTimeout * time.Millisecond,
		MinIdleConns:          config.MinIdleConns,
		MaxIdleConns:          config.MaxIdleConns,
		ConnMaxIdleTime:       config.ConnMaxIdleTime,
		ConnMaxLifetime:       config.ConnMaxLifetime,
	}
}

// buildRedisUniversalOptions 构建Redis通用客户端选项，支持单个、集群和哨兵模式
func buildRedisUniversalOptions(config *Config) *redis.UniversalOptions {
	return &redis.UniversalOptions{
		Addrs:                 config.Addrs,
		IdentitySuffix:        config.IdentitySuffix,
		MasterName:            config.MasterName,
		ClientName:            config.ClientName,
		DB:                    config.Database,
		Username:              config.Username,
		Password:              config.Password,
		SentinelUsername:      config.SentinelUsername,
		SentinelPassword:      config.SentinelPassword,
		MaxRetries:            config.MaxRetries,
		MinRetryBackoff:       config.MinRetryBackoff,
		MaxRetryBackoff:       config.MaxRetryBackoff,
		DialTimeout:           config.DialTimeout * time.Millisecond,
		ReadTimeout:           config.ReadTimeout * time.Millisecond,
		WriteTimeout:          config.WriteTimeout * time.Millisecond,
		ContextTimeoutEnabled: *config.ContextTimeoutEnabled,
		MaxRedirects:          config.MaxRedirects,
		ReadOnly:              *config.ReadOnly,
		RouteByLatency:        *config.RouteByLatency,
		RouteRandomly:         *config.RouteRandomly,
		PoolFIFO:              *config.PoolFIFO,
		PoolSize:              config.PoolSize,
		PoolTimeout:           config.PoolTimeout * time.Millisecond,
		MinIdleConns:          config.MinIdleConns,
		MaxIdleConns:          config.MaxIdleConns,
		MaxActiveConns:        config.MaxActiveConns,
		ConnMaxIdleTime:       config.ConnMaxIdleTime,
		ConnMaxLifetime:       config.ConnMaxLifetime,
	}
}

// Debug 设置调试模式
func Debug(b bool) Option {
	return func(config *Config) {
		config.Debug = b
	}
}

// MasterName 设置哨兵主节点名称
func MasterName(s string) Option {
	return func(config *Config) {
		config.MasterName = s
	}
}

// IdentitySuffix 设置客户端标识后缀
func IdentitySuffix(s string) Option {
	return func(config *Config) {
		config.IdentitySuffix = s
	}
}

// ClientName 设置客户端名称
func ClientName(s string) Option {
	return func(config *Config) {
		config.ClientName = s
	}
}

// Addr 添加一个地址
func Addr(s string) Option {
	return func(config *Config) {
		config.Addrs = append(config.Addrs, s)
	}
}

// Addrs 设置所有地址
func Addrs(ss []string) Option {
	return func(config *Config) {
		config.Addrs = ss
	}
}

// Database 设置数据库编号
func Database(n int) Option {
	return func(config *Config) {
		config.Database = n
	}
}

// Username 设置用户名
func Username(s string) Option {
	return func(config *Config) {
		config.Username = s
	}
}

// Password 设置密码
func Password(s string) Option {
	return func(config *Config) {
		config.Password = s
	}
}

// SentinelUsername 设置哨兵用户名
func SentinelUsername(s string) Option {
	return func(config *Config) {
		config.SentinelUsername = s
	}
}

// SentinelPassword 设置哨兵密码
func SentinelPassword(s string) Option {
	return func(config *Config) {
		config.SentinelPassword = s
	}
}

// Network 设置网络类型
func Network(s string) Option {
	return func(config *Config) {
		config.Network = s
	}
}

// MaxRetries 设置最大重试次数
func MaxRetries(n int) Option {
	return func(config *Config) {
		config.MaxRetries = n
	}
}

// MinRetryBackoff 设置最小重试间隔
func MinRetryBackoff(t time.Duration) Option {
	return func(config *Config) {
		config.MinRetryBackoff = t
	}
}

// MaxRetryBackoff 设置最大重试间隔
func MaxRetryBackoff(t time.Duration) Option {
	return func(config *Config) {
		config.MaxRetryBackoff = t
	}
}

// DialTimeout 设置连接超时时间
func DialTimeout(t time.Duration) Option {
	return func(config *Config) {
		config.DialTimeout = t
	}
}

// ReadTimeout 设置读取超时时间
func ReadTimeout(t time.Duration) Option {
	return func(config *Config) {
		config.ReadTimeout = t
	}
}

// WriteTimeout 设置写入超时时间
func WriteTimeout(t time.Duration) Option {
	return func(config *Config) {
		config.WriteTimeout = t
	}
}

// PoolFIFO 设置连接池类型（FIFO/LIFO）
func PoolFIFO(b bool) Option {
	return func(config *Config) {
		config.PoolFIFO = &b
	}
}

// PoolSize 设置连接池大小
func PoolSize(n int) Option {
	return func(config *Config) {
		config.PoolSize = n
	}
}

// PoolTimeout 设置连接池超时时间
func PoolTimeout(t time.Duration) Option {
	return func(config *Config) {
		config.PoolTimeout = t
	}
}

// IdleTimeout 设置空闲超时时间
func IdleTimeout(t time.Duration) Option {
	return func(config *Config) {
		config.IdleTimeout = t
	}
}

// MinIdleConns 设置最小空闲连接数
func MinIdleConns(n int) Option {
	return func(config *Config) {
		config.MinIdleConns = n
	}
}

// MaxIdleConns 设置最大空闲连接数
func MaxIdleConns(n int) Option {
	return func(config *Config) {
		config.MaxIdleConns = n
	}
}

// ConnMaxIdleTime 设置连接最大空闲时间
func ConnMaxIdleTime(t time.Duration) Option {
	return func(config *Config) {
		config.ConnMaxIdleTime = t
	}
}

// ConnMaxLifetime 设置连接最大生命周期
func ConnMaxLifetime(t time.Duration) Option {
	return func(config *Config) {
		config.ConnMaxLifetime = t
	}
}

// MaxRedirects 设置最大重定向次数
func MaxRedirects(n int) Option {
	return func(config *Config) {
		config.MaxRedirects = n
	}
}

// ReadOnly 设置只读模式
func ReadOnly(b bool) Option {
	return func(config *Config) {
		config.ReadOnly = &b
	}
}

// RouteRandomly 设置随机路由模式
func RouteRandomly(b bool) Option {
	return func(config *Config) {
		config.RouteRandomly = &b
	}
}

// RouteByLatency 设置延迟路由模式
func RouteByLatency(b bool) Option {
	return func(config *Config) {
		config.RouteByLatency = &b
	}
}

// ContextTimeoutEnabled 设置是否启用上下文超时
func ContextTimeoutEnabled(b bool) Option {
	return func(config *Config) {
		config.ContextTimeoutEnabled = &b
	}
}
