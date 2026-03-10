package resty

import (
	"context"

	"github.com/mel0dys0ng/song/internal/core/clients/resty"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

type (
	Client = resty.Client
	Option = resty.Option
)

// New return system http client based on resty.
// @key string redis config key.
func New(ctx context.Context, key string) *resty.Client {
	ctx = erlogs.StartTracef(ctx, "NewRestyClient:%s", key)
	defer erlogs.EndTrace(ctx, nil)
	return resty.CreateClient(ctx, "", key)
}

// Custom return the custom http client based on resty.
// @key string redis config key.
// @name string 自定义配置名称，全局唯一，否则后者覆盖前者。
// @opts []Option 自定义配置选项。
func Custom(ctx context.Context, name, key string, options ...resty.Option) *resty.Client {
	ctx = erlogs.StartTracef(ctx, "CustomRestyClient:%s:%s", name, key)
	defer erlogs.EndTrace(ctx, nil)
	return resty.CreateClient(ctx, name, key, options...)
}
