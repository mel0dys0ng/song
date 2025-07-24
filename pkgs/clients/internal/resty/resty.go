package resty

import (
	"context"
	"strings"

	resty2 "github.com/go-resty/resty/v2"
	"github.com/mel0dys0ng/song/pkgs/erlogs"
	"go.uber.org/zap"
)

// CreateClient
// @key string redis config key
// @name string 自定义配置名称。全局唯一，否则后者覆盖前者
// @opts []Option 自定义配置选项
func CreateClient(ctx context.Context, name, key string, opts ...Option) *Client {
	mk := key
	if len(name) > 0 {
		mk = strings.Join([]string{key, name}, "-")
	}

	if v, ok := clients.Load(mk); ok {
		if client, ok := v.(*Client); ok {
			return client
		}
	}

	elgSys := erlogs.New(erlogs.TypeSystem(), erlogs.Log(true), erlogs.Msgf("[resty] %s"))
	config, err := newConfig(ctx, key, elgSys, opts)
	if err != nil {
		err.PanicL(ctx, erlogs.Fields(zap.String("key", key)))
	}

	client := &Client{
		key:    key,
		config: config,
		Client: resty2.New(),
	}

	clients.Store(mk, client)

	return client
}
