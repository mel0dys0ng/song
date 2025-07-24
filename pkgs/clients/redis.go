package clients

import (
	"context"
	"fmt"

	"github.com/mel0dys0ng/song/pkgs/clients/internal/redis"
	"github.com/mel0dys0ng/song/pkgs/erlogs"
)

type (
	RedisClient          = redis.Client
	RedisUniversalClient = redis.UniversalClient
)

// NewRedisClient return the normal redis client
// @key string redis config key
func NewRedisClient(ctx context.Context, key string) *redis.Client {
	ctx = erlogs.StartTrace(ctx, fmt.Sprintf("NewRedisClient:%s", key))
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateClient(ctx, "", key)
}

// NewRedisUniversalClient return the universal redis client
// @key string redis config key
func NewRedisUniversalClient(ctx context.Context, key string) *redis.UniversalClient {
	ctx = erlogs.StartTrace(ctx, fmt.Sprintf("NewRedisUniversalClient:%s", key))
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateUniversalClient(ctx, "", key)
}

// CustomRedisClient return the normal redis client with custom options
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @options []Option 自定义配置选项
func CustomRedisClient(ctx context.Context, name, key string, options ...redis.Option) *redis.Client {
	ctx = erlogs.StartTrace(ctx, fmt.Sprintf("CustomRedisClient:%s:%s", name, key))
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateClient(ctx, name, key, options...)
}

// CustomRedisUniversalClient return the universal redis client with custom options
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @options []Option 自定义配置选项
func CustomRedisUniversalClient(ctx context.Context, name, key string, options ...redis.Option) *redis.UniversalClient {
	ctx = erlogs.StartTrace(ctx, fmt.Sprintf("CustomRedisUniversalClient:%s:%s", name, key))
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateUniversalClient(ctx, name, key, options...)
}
