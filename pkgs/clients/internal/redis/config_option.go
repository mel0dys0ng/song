package redis

import (
	"time"

	"github.com/mel0dys0ng/song/pkgs/utils/sljces"
	"github.com/redis/go-redis/v9"
)

const (
	NetWorkTCP                   = "tcp"
	NetworkUnix                  = "unix"
	DefauktDebug                 = false
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
	DefaultConnMaxIdleTime       = 1000
	DefaultConnMaxLifetime       = 1000
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

type Option struct {
	Func func(*Config)
}

func buildRedisOptions(config *Config) *redis.Options {
	return &redis.Options{
		Network:               config.Network,
		Addr:                  sljces.First(config.Addrs, ""),
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
		ContextTimeoutEnabled: config.ContextTimeoutEnabled,
		PoolFIFO:              config.PoolFIFO,
		PoolSize:              config.PoolSize,
		PoolTimeout:           config.PoolTimeout * time.Millisecond,
		MinIdleConns:          config.MinIdleConns,
		MaxIdleConns:          config.MaxIdleConns,
		ConnMaxIdleTime:       config.ConnMaxIdleTime,
		ConnMaxLifetime:       config.ConnMaxLifetime,
	}
}

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
		ContextTimeoutEnabled: config.ContextTimeoutEnabled,
		MaxRedirects:          config.MaxRedirects,
		ReadOnly:              config.ReadOnly,
		RouteByLatency:        config.RouteByLatency,
		RouteRandomly:         config.RouteRandomly,
		PoolFIFO:              config.PoolFIFO,
		PoolSize:              config.PoolSize,
		PoolTimeout:           config.PoolTimeout * time.Millisecond,
		MinIdleConns:          config.MinIdleConns,
		MaxIdleConns:          config.MaxIdleConns,
		MaxActiveConns:        config.MaxActiveConns,
		ConnMaxIdleTime:       config.ConnMaxIdleTime,
		ConnMaxLifetime:       config.ConnMaxLifetime,
	}
}

func Debug(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.Debug = b
		},
	}
}

func MasterName(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.MasterName = s
		},
	}
}

func IdentitySuffix(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.IdentitySuffix = s
		},
	}
}

func ClientName(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.ClientName = s
		},
	}
}

func Addr(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.Addrs = append(config.Addrs, s)
		},
	}
}

func Addrs(ss []string) Option {
	return Option{
		Func: func(config *Config) {
			config.Addrs = ss
		},
	}
}

func Database(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.Database = n
		},
	}
}

func Username(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.Username = s
		},
	}
}

func Password(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.Password = s
		},
	}
}

func SentinelUsername(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.SentinelUsername = s
		},
	}
}

func SentinelPassword(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.SentinelPassword = s
		},
	}
}

func Network(s string) Option {
	return Option{
		Func: func(config *Config) {
			config.Network = s
		},
	}
}

func MaxRetries(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxRetries = n
		},
	}
}

func MinRetryBackoff(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.MinRetryBackoff = t
		},
	}
}

func MaxRetryBackoff(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxRetryBackoff = t
		},
	}
}

func DialTimeout(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.DialTimeout = t
		},
	}
}

func ReadTimeout(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.ReadTimeout = t
		},
	}
}

func WriteTimeout(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.WriteTimeout = t
		},
	}
}

func PoolFIFO(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.PoolFIFO = b
		},
	}
}

func PoolSize(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.PoolSize = n
		},
	}
}

func PoolTimeout(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.PoolTimeout = t
		},
	}
}

func IdleTimeout(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.IdleTimeout = t
		},
	}
}

func MinIdleConns(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.MinIdleConns = n
		},
	}
}

func MaxIdleConns(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxIdleConns = n
		},
	}
}

func ConnMaxIdleTime(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.ConnMaxIdleTime = t
		},
	}
}

func ConnMaxLifetime(t time.Duration) Option {
	return Option{
		Func: func(config *Config) {
			config.ConnMaxLifetime = t
		},
	}
}

func MaxRedirects(n int) Option {
	return Option{
		Func: func(config *Config) {
			config.MaxRedirects = n
		},
	}
}

func ReadOnly(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.ReadOnly = b
		},
	}
}

func RouteRandomly(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.RouteRandomly = b
		},
	}
}

func RouteByLatency(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.RouteByLatency = b
		},
	}
}

func ContextTimeoutEnabled(b bool) Option {
	return Option{
		Func: func(config *Config) {
			config.ContextTimeoutEnabled = b
		},
	}
}
