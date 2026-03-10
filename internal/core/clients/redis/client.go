package redis

import (
	"context"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/singleton"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Client Redis客户端结构体
type Client struct {
	// 普通模式redis client
	*redis.Client
	// redis config key
	key string
	// redis config config
	config *Config
}

// Key 返回Redis配置的键
func (c *Client) Key() string {
	return c.key
}

// CreateClient 创建 Redis 客户端实例
// @param ctx 上下文
// @param name 自定义配置名称，全局唯一，否则后者覆盖前者
// @param key Redis 配置键
// @param opts 可选配置参数
// @return *Client Redis 客户端实例
func CreateClient(ctx context.Context, name, key string, opts ...Option) *Client {
	return singleton.Once(singleton.Key(), func() *Client {
		config, err := NewConfig(ctx, key, opts)
		if err != nil {
			erlogs.Convert(err).Wrap("failed to load options: UnmarshalKey error").Options(BaseELOptions()).PanicLog(ctx,
				erlogs.OptionFields(zap.String("key", key)),
			)
		}

		client := &Client{
			key:    key,
			config: config,
			Client: redis.NewClient(
				buildRedisOptions(config),
			),
		}

		return client
	})
}
