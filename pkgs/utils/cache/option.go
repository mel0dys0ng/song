package cache

import (
	"time"

	"github.com/mel0dys0ng/song/pkgs/utils/retry"
	"github.com/redis/go-redis/v9"
)

func RedisCache[T any](client redis.UniversalClient, ttl time.Duration) Option[T] {
	return Option[T]{
		apply: func(c *Cache[T]) {
			c.redisCache = newRedisCache[T](client, ttl)
		},
	}
}

func LRUCache[T any](size int, ttl time.Duration) Option[T] {
	return Option[T]{
		apply: func(c *Cache[T]) {
			c.lruCache = newLRUCache[T](size, ttl)
		},
	}
}

func Retry[T any](enable, singleflight bool, opts ...retry.Option) Option[T] {
	return Option[T]{
		apply: func(c *Cache[T]) {
			c.retryConf = retryConf{
				enable:       enable,
				singleflight: singleflight,
				options:      opts,
			}
		},
	}
}

func IsZero[T any](fn func(data T) bool) Option[T] {
	return Option[T]{
		apply: func(c *Cache[T]) {
			c.isZero = fn
		},
	}
}

func DataId[T any](fn func(data T) any) Option[T] {
	return Option[T]{
		apply: func(c *Cache[T]) {
			c.dataId = fn
		},
	}
}

func KeyPrefix[T any](s string) Option[T] {
	return Option[T]{
		apply: func(c *Cache[T]) {
			c.keyPrefix = s
		},
	}
}
