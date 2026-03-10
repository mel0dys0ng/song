package redis

import (
	"time"

	"github.com/mel0dys0ng/song/internal/core/clients/redis"
)

func OptionDebug(b bool) redis.Option {
	return redis.Debug(b)
}

func OptionClientName(s string) redis.Option {
	return redis.ClientName(s)
}

func OptionAddr(s string) redis.Option {
	return redis.Addr(s)
}

func OptionAddrs(ss []string) redis.Option {
	return redis.Addrs(ss)
}

func OptionDatabase(n int) redis.Option {
	return redis.Database(n)
}

func OptionUsername(s string) redis.Option {
	return redis.Username(s)
}

func OptionPassword(s string) redis.Option {
	return redis.Password(s)
}

func OptionNetwork(s string) redis.Option {
	return redis.Network(s)
}

func OptionMaxRetries(n int) redis.Option {
	return redis.MaxRetries(n)
}

func OptionMinRetryBackoff(t time.Duration) redis.Option {
	return redis.MinRetryBackoff(t)
}

func OptionMaxRetryBackoff(t time.Duration) redis.Option {
	return redis.MaxRetryBackoff(t)
}

func OptionDialTimeout(t time.Duration) redis.Option {
	return redis.DialTimeout(t)
}

func OptionReadTimeout(t time.Duration) redis.Option {
	return redis.ReadTimeout(t)
}

func OptionWriteTimeout(t time.Duration) redis.Option {
	return redis.WriteTimeout(t)
}

func OptionPoolFIFO(b bool) redis.Option {
	return redis.PoolFIFO(b)
}

func OptionPoolSize(n int) redis.Option {
	return redis.PoolSize(n)
}

func OptionPoolTimeout(t time.Duration) redis.Option {
	return redis.PoolTimeout(t)
}

func OptionIdleTimeout(t time.Duration) redis.Option {
	return redis.IdleTimeout(t)
}

func OptionMinIdleConns(n int) redis.Option {
	return redis.MinIdleConns(n)
}

func OptionMaxIdleConns(n int) redis.Option {
	return redis.MaxIdleConns(n)
}

func OptionConnMaxIdleTime(t time.Duration) redis.Option {
	return redis.ConnMaxIdleTime(t)
}

func OptionConnMaxLifetime(t time.Duration) redis.Option {
	return redis.ConnMaxLifetime(t)
}

func OptionMaxRedirects(n int) redis.Option {
	return redis.MaxRedirects(n)
}

func OptionReadOnly(b bool) redis.Option {
	return redis.ReadOnly(b)
}

func OptionRouteRandomly(b bool) redis.Option {
	return redis.RouteRandomly(b)
}

func OptionRouteByLatency(b bool) redis.Option {
	return redis.RouteByLatency(b)
}

func OptionContextTimeoutEnabled(b bool) redis.Option {
	return redis.ContextTimeoutEnabled(b)
}

func OptionIdentitySuffix(s string) redis.Option {
	return redis.IdentitySuffix(s)
}

func OptionMasterName(s string) redis.Option {
	return redis.MasterName(s)
}

func OptionSentinelUsername(s string) redis.Option {
	return redis.SentinelUsername(s)
}

func OptionSentinelPassword(s string) redis.Option {
	return redis.SentinelPassword(s)
}
