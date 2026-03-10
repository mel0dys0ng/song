package redis

import (
	"context"

	"github.com/mel0dys0ng/song/internal/core/clients/redis"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

type (
	Client          = redis.Client
	UniversalClient = redis.UniversalClient
	Option          = redis.Option
)

// NewClient return the normal redis client
// @key string redis config key
func NewClient(ctx context.Context, key string) *redis.Client {
	ctx = erlogs.StartTracef(ctx, "NewRedisClient:%s", key)
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateClient(ctx, "", key)
}

// NewUniversalClient return the universal redis client
// @key string redis config key
func NewUniversalClient(ctx context.Context, key string) *redis.UniversalClient {
	ctx = erlogs.StartTracef(ctx, "NewRedisUniversalClient:%s", key)
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateUniversalClient(ctx, "", key)
}

// CustomClient return the normal redis client with custom options
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @options []Option 自定义配置选项
func CustomClient(ctx context.Context, name, key string, options ...redis.Option) *redis.Client {
	ctx = erlogs.StartTracef(ctx, "CustomRedisClient:%s:%s", name, key)
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateClient(ctx, name, key, options...)
}

// CustomUniversalClient return the universal redis client with custom options
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @options []Option 自定义配置选项
func CustomUniversalClient(ctx context.Context, name, key string, options ...redis.Option) *redis.UniversalClient {
	ctx = erlogs.StartTracef(ctx, "CustomRedisUniversalClient:%s:%s", name, key)
	defer erlogs.EndTrace(ctx, nil)
	return redis.CreateUniversalClient(ctx, name, key, options...)
}
