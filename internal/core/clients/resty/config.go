package resty

import (
	"context"
	"time"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/mel0dys0ng/song/pkg/vipers"
)

// Config 定义 Resty 客户端的配置
type Config struct {
	// Debug 是否开启调试模式
	Debug bool `json:"debug" yaml:"debug" mapstructure:"debug"`

	// Type 网络类型，extranet or intranet
	Type string `json:"type" yaml:"type" mapstructure:"type"`

	// Trace 是否开启追踪
	Trace bool `json:"trace" yaml:"trace" mapstructure:"trace"`

	// BaseURL 请求的基础URL
	BaseURL string `json:"baseUrl" yaml:"baseUrl" mapstructure:"baseUrl"`

	// Did 依赖服务ID
	Did string `json:"did" yaml:"did" mapstructure:"did"`

	// Timeout 请求超时时间，默认500毫秒
	Timeout time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`

	// RetryCount 重试次数，默认3次
	RetryCount int `json:"retryCount" yaml:"retryCount" mapstructure:"retryCount"`

	// RetryWaitTime 重试等待时间，默认100毫秒
	RetryWaitTime time.Duration `json:"retryWaitTime" yaml:"retryWaitTime" mapstructure:"retryWaitTime"`

	// RetryWaitMaxTime 重试最大等待时间，默认2秒
	RetryWaitMaxTime time.Duration `json:"retryWaitMaxTime" yaml:"retryWaitMaxTime" mapstructure:"retryWaitMaxTime"`

	// SignSecret 请求签名密钥
	SignSecret string `json:"signSecret" yaml:"signSecret" mapstructure:"signSecret"`

	// SignTTL 请求签名有效期，单位秒
	SignTTL int `json:"signTTL" yaml:"signTTL" mapstructure:"signTTL"`

	// SignConfig 签名配置
	SignConfig *SignConfig `json:"signConfig" yaml:"signConfig" mapstructure:"signConfig"`
}

// SignConfig 签名配置结构体
type SignConfig struct {
	Enable   bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
	Secret   string `json:"secret" yaml:"secret" mapstructure:"secret"`
	TTL      int    `json:"ttl" yaml:"ttl" mapstructure:"ttl"` // TTL in seconds
	Query    bool   `json:"query" yaml:"query" mapstructure:"query"`
	FormData bool   `json:"formData" yaml:"formData" mapstructure:"formData"`
	Header   bool   `json:"header" yaml:"header" mapstructure:"header"`
}

// newConfig 创建新的配置实例
func newConfig(ctx context.Context, key string, opts []Option) (config *Config, err error) {
	// 初始化默认配置
	config = &Config{
		Debug:            DefaultDebug,
		Did:              DefaultDid,
		Type:             DefaultType,
		Trace:            DefaultTrace,
		Timeout:          DefaultTimeout,
		RetryCount:       DefaultRetryCount,
		RetryWaitTime:    DefaultRetryWaitTime,
		RetryWaitMaxTime: DefaultRetryWaitMaxTime,
		SignSecret:       DefaultSignSecret,
		SignTTL:          DefaultSignTTL,
		SignConfig: &SignConfig{
			Enable:   false,
			Secret:   DefaultSignSecret,
			TTL:      DefaultSignTTL,
			Query:    true,
			FormData: true,
			Header:   true,
		},
	}

	// 加载配置选项
	options, err := config.LoadOptions(ctx, key)
	if err != nil {
		return
	}

	// 合并默认选项和自定义选项
	opts = append(options, opts...)
	for _, v := range opts {
		v(config)
	}

	// 验证配置
	err = config.check()

	return
}

// LoadOptions 从配置源加载配置选项
func (c *Config) LoadOptions(ctx context.Context, key string) (opts []Option, err error) {

	// 创建配置值容器
	values := &Config{}
	er := vipers.UnmarshalKey(key, values)
	if er != nil {
		err = erlogs.Convert(er).Wrap("failed to unmarshal key").Panic()
		return
	}

	// 根据配置值生成选项
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

	if values.SignSecret != DefaultSignSecret {
		opts = append(opts, SignSecret(values.SignSecret))
	}

	return
}

// check 验证配置的有效性
func (c *Config) check() (err error) {
	// 验证基础 URL 是否为空
	if len(c.BaseURL) == 0 {
		err = erlogs.New("基础 URL 不能为空").Panic()
		return
	}

	// 如果不是外网请求，则设置为内网请求
	if c.Type != Extranet {
		c.Type = Intranet
	}

	return
}
