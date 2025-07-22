package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/mel0dys0ng/song/erlogs"
	"github.com/mel0dys0ng/song/utils/singleton"
	"github.com/mel0dys0ng/song/utils/strjngs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	_elgSys erlogs.ErLogInterface
	once    sync.Once
)

type Client struct {
	// 普通模式redis client
	*redis.Client
	// redis config key
	key string
	// redis config config
	config *Config
}

// Key return the normal redis config key
func (c *Client) Key() string {
	return c.key
}

func elgSys() erlogs.ErLogInterface {
	once.Do(func() {
		_elgSys = erlogs.New(erlogs.TypeSystem(), erlogs.Log(true), erlogs.Msgf("[redis] %s"))
	})
	return _elgSys
}

// CreateClient
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @opts []Option 自定义配置选项
func CreateClient(ctx context.Context, name, key string, opts ...Option) *Client {
	instanceKey := fmt.Sprintf("redis-client-%s", strjngs.GenerateStableUniqueStr(name, key))
	return singleton.GetInstance(instanceKey, func() *Client {
		config, err := NewConfig(ctx, key, elgSys(), opts)
		if err != nil {
			err.PanicL(ctx, erlogs.Fields(zap.String("key", key)))
		}

		return &Client{
			key:    key,
			config: config,
			Client: redis.NewClient(
				buildRedisOptions(config),
			),
		}
	})
}
