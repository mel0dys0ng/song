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
	DefaultSignTTL          = 10 * time.Second
	DefaultSignSecret       = ""
)

type Option struct {
	Apply func(*Config)
}

func Debug(b bool) Option {
	return Option{func(c *Config) { c.Debug = b }}
}

func Type(s string) Option {
	return Option{
		Apply: func(c *Config) {
			c.Type = s
		},
	}
}

func Did(s string) Option {
	return Option{
		Apply: func(c *Config) {
			c.Did = s
		},
	}
}

func BaseURL(s string) Option {
	return Option{
		Apply: func(c *Config) {
			c.BaseURL = s
		},
	}
}

func Trace(b bool) Option {
	return Option{
		Apply: func(c *Config) {
			c.Trace = b
		},
	}
}

func Timeout(t time.Duration) Option {
	return Option{
		Apply: func(c *Config) {
			c.Timeout = t
		},
	}
}

func RetryCount(i int) Option {
	return Option{
		Apply: func(c *Config) {
			c.RetryCount = i
		},
	}
}

func RetryWaitTime(t time.Duration) Option {
	return Option{
		Apply: func(c *Config) {
			c.RetryWaitTime = t
		},
	}
}

func RetryWaitMaxTime(t time.Duration) Option {
	return Option{
		Apply: func(c *Config) {
			c.RetryWaitMaxTime = t
		},
	}
}

func SignSecret(s string) Option {
	return Option{
		Apply: func(c *Config) {
			c.SignSerect = s
		},
	}
}

func SignTTL(t time.Duration) Option {
	return Option{
		Apply: func(c *Config) {
			c.SignTTL = t
		},
	}
}
