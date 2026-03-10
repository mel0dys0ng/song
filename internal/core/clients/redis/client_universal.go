package redis

import (
	"context"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/singleton"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// UniversalClient 通用Redis客户端，支持单个、集群和哨兵模式
type UniversalClient struct {
	// 通用redis client，支持单个、集群和哨兵模式的redis
	redis.UniversalClient
	// redis config key
	key string
	// redis config config
	config *Config
}

// Key 返回通用Redis配置的键
func (c *UniversalClient) Key() string {
	return c.key
}

// CreateUniversalClient 创建通用 Redis 客户端实例
// @param ctx 上下文
// @param name 自定义配置名称，全局唯一，否则后者覆盖前者
// @param key Redis 配置键
// @param opts 可选配置参数
// @return *UniversalClient 通用 Redis 客户端实例
func CreateUniversalClient(ctx context.Context, name, key string, opts ...Option) *UniversalClient {
	return singleton.Once(singleton.Key(), func() *UniversalClient {
		config, err := NewConfig(ctx, key, opts)
		if err != nil {
			erlogs.Convert(err).Wrap("failed to create config").Options(BaseELOptions()).PanicLog(ctx,
				erlogs.OptionFields(zap.String("key", key)),
			)
		}

		universalClient := &UniversalClient{
			key:    key,
			config: config,
			UniversalClient: redis.NewUniversalClient(
				buildRedisUniversalOptions(config),
			),
		}

		return universalClient
	})
}
