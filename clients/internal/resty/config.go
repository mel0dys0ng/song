package resty

import (
	"context"
	"time"

	"github.com/mel0dys0ng/song/erlogs"
	"github.com/mel0dys0ng/song/vipers"
)

type Config struct {
	// Debug 是否开启调试
	Debug bool `json:"debug" mapstructure:"debug"`

	// Type 网络类型，extranet or intranet
	Type string `json:"type" mapstructure:"type"`

	// Trace 是否开启trace
	Trace bool `json:"trace" mapstructure:"trace"`

	// BaseURL request base url
	BaseURL string `json:"baseUrl" mapstructure:"baseUrl"`

	// Did Dependency ID
	Did string `json:"did" mapstructure:"did"`

	// Timeout 超时时间，默认500毫秒
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`

	// RetryCount 重试次数，默认3次
	RetryCount int `json:"retryCount" mapstructure:"retryCount"`

	// RetryWaitTime 重试等待时间，默认100毫秒
	RetryWaitTime time.Duration `json:"retryWaitTime" mapstructure:"retryWaitTime"`

	// RetryWaitMaxTime 重试最大等待时间，默认2秒
	RetryWaitMaxTime time.Duration `json:"retryWaitMaxTime" mapstructure:"retryWaitMaxTime"`

	// SignSerect 请求签名secret
	SignSerect string `json:"signSerect" mapstructure:"signSerect"`

	// SignTTL 请求签名TTL
	SignTTL time.Duration `json:"signTTL" mapstructure:"signTTL"`
}

func newConfig(ctx context.Context, key string, elg erlogs.ErLogInterface, opts []Option) (
	config *Config, err erlogs.ErLogInterface) {

	config = &Config{
		Debug:            DefaultDebug,
		Did:              DefaultDid,
		Type:             DefaultType,
		Trace:            DefaultTrace,
		Timeout:          DefaultTimeout,
		RetryCount:       DefaultRetryCount,
		RetryWaitTime:    DefaultRetryWaitTime,
		RetryWaitMaxTime: DefaultRetryWaitMaxTime,
		SignSerect:       DefaultSignSecret,
		SignTTL:          DefaultSignTTL,
	}

	options, err := config.LoadOptions(ctx, key, elg)
	if err != nil {
		return
	}

	opts = append(options, opts...)
	for _, v := range opts {
		v.Apply(config)
	}

	err = config.check(ctx, elg)

	return
}

// LoadOptions 加载配置并返回Options
func (c *Config) LoadOptions(ctx context.Context, key string, elg erlogs.ErLogInterface) (
	opts []Option, err erlogs.ErLogInterface) {

	values := &Config{}
	er := vipers.UnmarshalKey(key, values, nil)
	if er != nil {
		err = elg.PanicE(ctx,
			erlogs.Msgv("failed to load options: UnmarshalKey error"),
			erlogs.Content(er.Error()),
		)
		return
	}

	if values.Debug != DefaultDebug {
		opts = append(opts, Debug(values.Debug))
	}

	if values.Did != DefaultDid {
		opts = append(opts, Did(values.Did))
	}

	if values.BaseURL != DefaultBaseURL {
		opts = append(opts, BaseURL(values.BaseURL))
	}

	if values.RetryWaitMaxTime != DefaultRetryWaitMaxTime {
		opts = append(opts, RetryWaitMaxTime(values.RetryWaitMaxTime))
	}

	if values.RetryWaitTime != DefaultRetryWaitTime {
		opts = append(opts, RetryWaitTime(values.RetryWaitTime))
	}

	if values.RetryCount != DefaultRetryCount {
		opts = append(opts, RetryCount(values.RetryCount))
	}

	if values.Timeout != DefaultTimeout {
		opts = append(opts, Timeout(values.Timeout))
	}

	if values.Type != DefaultType {
		opts = append(opts, Type(values.Type))
	}

	if values.SignTTL != DefaultSignTTL {
		opts = append(opts, SignTTL(values.SignTTL))
	}

	if values.SignSerect != DefaultSignSecret {
		opts = append(opts, SignSecret(values.SignSerect))
	}

	return
}

func (c *Config) check(ctx context.Context, elg erlogs.ErLogInterface) (err erlogs.ErLogInterface) {
	if len(c.BaseURL) == 0 {
		err = elg.PanicE(ctx,
			erlogs.Msgv("baseURL should not be empty"),
		)
		return
	}

	if c.Type != Extranet {
		c.Type = Intranet
	}

	return
}
