package cache

import (
	"time"

	"github.com/mel0dys0ng/song/pkg/retry"
	"github.com/redis/go-redis/v9"
)

// RedisCache 配置Redis缓存选项
func RedisCache[T any](client redis.UniversalClient, ttl time.Duration) Option[T] {
	return func(c *Cache[T]) {
		c.redisCache = newRedisCache[T](client, ttl)
	}
}

// LRUCache 配置LRU缓存选项
func LRUCache[T any](size int, ttl time.Duration) Option[T] {
	return func(c *Cache[T]) {
		c.lruCache = newLRUCache[T](size, ttl)
	}
}

// Retry 配置重试选项
func Retry[T any](enable, singleflight bool, opts ...retry.Option) Option[T] {
	return func(c *Cache[T]) {
		c.retryConf = retryConf{
			enable:       enable,
			singleflight: singleflight,
			options:      opts,
		}
	}
}

// IsZero 配置零值检查函数
func IsZero[T any](fn func(data T) bool) Option[T] {
	return func(c *Cache[T]) {
		c.isZero = fn
	}
}

// DataId 配置数据ID提取函数
func DataId[T any](fn func(data T) any) Option[T] {
	return func(c *Cache[T]) {
		c.dataId = fn
	}
}

// KeyPrefix 配置键前缀
func KeyPrefix[T any](s string) Option[T] {
	return func(c *Cache[T]) {
		c.keyPrefix = s
	}
}
