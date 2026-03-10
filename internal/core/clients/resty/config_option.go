package resty

import (
	"time"
)

const (
	Intranet                = "intranet" // 内网请求，默认
	Extranet                = "extranet" // 外网请求
	DefaultDebug            = false
	DefaultBaseURL          = ""
	DefaultDid              = ""
	DefaultType             = Intranet
	DefaultTrace            = true
	DefaultRetryCount       = 3
	DefaultTimeout          = 500 * time.Millisecond
	DefaultRetryWaitTime    = 100 * time.Millisecond
	DefaultRetryWaitMaxTime = 2 * time.Second
	DefaultSignTTL          = 300                   // 5分钟，与https模块保持一致
	DefaultSignSecret       = "default_sign_secret" // 与https模块保持一致
)

// Option 用于配置 Resty 客户端选项的函数类型
type Option func(*Config)

// Debug 设置调试模式
func Debug(b bool) Option {
	return func(c *Config) { c.Debug = b }
}

// Type 设置网络类型
func Type(s string) Option {
	return func(c *Config) {
		c.Type = s
	}
}

// Did 设置依赖服务ID
func Did(s string) Option {
	return func(c *Config) {
		c.Did = s
	}
}

// BaseURL 设置请求的基础URL
func BaseURL(s string) Option {
	return func(c *Config) {
		c.BaseURL = s
	}
}

// Trace 设置是否开启追踪
func Trace(b bool) Option {
	return func(c *Config) {
		c.Trace = b
	}
}

// Timeout 设置请求超时时间
func Timeout(t time.Duration) Option {
	return func(c *Config) {
		c.Timeout = t
	}
}

// RetryCount 设置重试次数
func RetryCount(i int) Option {
	return func(c *Config) {
		c.RetryCount = i
	}
}

// RetryWaitTime 设置重试等待时间
func RetryWaitTime(t time.Duration) Option {
	return func(c *Config) {
		c.RetryWaitTime = t
	}
}

// RetryWaitMaxTime 设置重试最大等待时间
func RetryWaitMaxTime(t time.Duration) Option {
	return func(c *Config) {
		c.RetryWaitMaxTime = t
	}
}

// SignSecret 设置签名密钥
func SignSecret(s string) Option {
	return func(c *Config) {
		c.SignSecret = s
	}
}

// SignTTL 设置签名有效期（秒）
func SignTTL(t int) Option {
	return func(c *Config) {
		c.SignTTL = t
	}
}
