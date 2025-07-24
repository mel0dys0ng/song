package clients

import (
	"time"

	"github.com/mel0dys0ng/song/pkgs/clients/internal/redis"
)

func RedisOptionDebug(b bool) redis.Option {
	return redis.Debug(b)
}

func RedisOptionClientName(s string) redis.Option {
	return redis.ClientName(s)
}

func RedisOptionAddr(s string) redis.Option {
	return redis.Addr(s)
}

func RedisOptionAddrs(ss []string) redis.Option {
	return redis.Addrs(ss)
}

func RedisOptionDatabase(n int) redis.Option {
	return redis.Database(n)
}

func RedisOptionUsername(s string) redis.Option {
	return redis.Username(s)
}

func RedisOptionPassword(s string) redis.Option {
	return redis.Password(s)
}

func RedisOptionNetwork(s string) redis.Option {
	return redis.Network(s)
}

func RedisOptionMaxRetries(n int) redis.Option {
	return redis.MaxRetries(n)
}

func RedisOptionMinRetryBackoff(t time.Duration) redis.Option {
	return redis.MinRetryBackoff(t)
}

func RedisOptionMaxRetryBackoff(t time.Duration) redis.Option {
	return redis.MaxRetryBackoff(t)
}

func RedisOptionDialTimeout(t time.Duration) redis.Option {
	return redis.DialTimeout(t)
}

func RedisOptionReadTimeout(t time.Duration) redis.Option {
	return redis.ReadTimeout(t)
}

func RedisOptionWriteTimeout(t time.Duration) redis.Option {
	return redis.WriteTimeout(t)
}

func RedisOptionPoolFIFO(b bool) redis.Option {
	return redis.PoolFIFO(b)
}

func RedisOptionPoolSize(n int) redis.Option {
	return redis.PoolSize(n)
}

func RedisOptionPoolTimeout(t time.Duration) redis.Option {
	return redis.PoolTimeout(t)
}

func RedisOptionIdleTimeout(t time.Duration) redis.Option {
	return redis.IdleTimeout(t)
}

func RedisOptionMinIdleConns(n int) redis.Option {
	return redis.MinIdleConns(n)
}

func RedisOptionMaxIdleConns(n int) redis.Option {
	return redis.MaxIdleConns(n)
}

func RedisOptionConnMaxIdleTime(t time.Duration) redis.Option {
	return redis.ConnMaxIdleTime(t)
}

func RedisOptionConnMaxLifetime(t time.Duration) redis.Option {
	return redis.ConnMaxLifetime(t)
}

func RedisOptionMaxRedirects(n int) redis.Option {
	return redis.MaxRedirects(n)
}

func RedisOptionReadOnly(b bool) redis.Option {
	return redis.ReadOnly(b)
}

func RedisOptionRouteRandomly(b bool) redis.Option {
	return redis.RouteRandomly(b)
}

func RedisOptionRouteByLatency(b bool) redis.Option {
	return redis.RouteByLatency(b)
}

func RedisOptionContextTimeoutEnabled(b bool) redis.Option {
	return redis.ContextTimeoutEnabled(b)
}
