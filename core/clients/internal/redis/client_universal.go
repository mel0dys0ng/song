package redis

import (
	"context"
	"fmt"

	"github.com/mel0dys0ng/song/core/erlogs"
	"github.com/mel0dys0ng/song/core/utils/singleton"
	"github.com/mel0dys0ng/song/core/utils/strjngs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type UniversalClient struct {
	// 通用redis client，支持单个、集群和哨兵模式的redis
	redis.UniversalClient
	// redis config key
	key string
	// redis config config
	config *Config
}

// Key return the universal redis config key
func (c *UniversalClient) Key() string {
	return c.key
}

// CreateUniversalClient
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @opts []Option 自定义配置选项
func CreateUniversalClient(ctx context.Context, name, key string, opts ...Option) *UniversalClient {
	instanceKey := fmt.Sprintf("redis-universal-client-%s", strjngs.GenerateStableUniqueStr(name, key))
	return singleton.GetInstance(instanceKey, func() *UniversalClient {
		config, err := NewConfig(ctx, key, elgSys(), opts)
		if err != nil {
			err.PanicL(ctx, erlogs.Fields(zap.String("key", key)))
		}

		return &UniversalClient{
			key:    key,
			config: config,
			UniversalClient: redis.NewUniversalClient(
				buildRedisUniversalOptions(config),
			),
		}
	})
}
