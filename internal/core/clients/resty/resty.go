package resty

import (
	"context"
	"strings"

	resty2 "github.com/go-resty/resty/v2"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/metas"
	"go.uber.org/zap"
)

// CreateClient 创建 Resty 客户端实例
// 参数 name 为自定义配置名称，与key组合形成唯一标识
// 参数 key 为配置键，用于从配置源加载配置
// 参数 opts 为可选的配置选项
// 返回创建的客户端实例
func CreateClient(ctx context.Context, name, key string, opts ...Option) *Client {
	// 构造唯一标识键
	mk := key
	if len(name) > 0 {
		mk = strings.Join([]string{key, name}, "-")
	}

	// 尝试从缓存中获取已存在的客户端
	if v, ok := clients.Load(mk); ok {
		if client, ok := v.(*Client); ok {
			return client
		}
	}

	// 创建配置
	config, err := newConfig(ctx, key, opts)
	if err != nil {
		erlogs.Convert(err).Wrap("failed to create config").Options(BaseELOptions()).PanicLog(ctx,
			erlogs.OptionFields(zap.String("key", key)),
		)
	}

	// 创建 Resty 客户端实例
	client := &Client{
		key:      key,
		config:   config,
		Client:   resty2.New(),
		metadata: metas.Metadata(),
	}

	// 设置基本配置
	client.SetBaseURL(config.BaseURL)
	client.SetTimeout(config.Timeout)
	client.SetRetryCount(config.RetryCount)
	client.SetRetryWaitTime(config.RetryWaitTime)
	client.SetRetryMaxWaitTime(config.RetryWaitMaxTime)

	// 存储到缓存中
	clients.Store(mk, client)

	return client
}
