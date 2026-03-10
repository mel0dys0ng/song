package https

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/internal/core/https"
	"github.com/samber/lo"
)

// Port 设置http server端口
func Port(i int) Option {
	return func(options *https.Options) {
		options.Port = i
	}
}

// TLSOpen 是否开启https
func TLSOpen(b bool) Option {
	return func(options *https.Options) {
		options.TLSOpen = b
	}
}

// TLSKeyFile 设置https TLS key file path
func TLSKeyFile(s string) Option {
	return func(options *https.Options) {
		options.TLSKeyFile = s
	}
}

// TLSKeyCert 设置https TLS cert file path
func TLSKeyCert(s string) Option {
	return func(options *https.Options) {
		options.TLSCertFile = s
	}
}

func KeepAlive(b bool) Option {
	return func(options *https.Options) {
		options.KeepAlive = b
	}
}

func ReadTimeout(s string) Option {
	return func(options *https.Options) {
		options.ReadTimeout = s
	}
}

func ReadHeaderTimeout(s string) Option {
	return func(options *https.Options) {
		options.ReadHeaderTimeout = s
	}
}

func WriteTimeout(s string) Option {
	return func(options *https.Options) {
		options.WriteTimeout = s
	}
}

func IdleTimeout(s string) Option {
	return func(options *https.Options) {
		options.IdleTimeout = s
	}
}

func HammerTime(s string) Option {
	return func(options *https.Options) {
		options.HammerTime = s
	}
}

func LoggerHeaderKeys(keys ...string) Option {
	return func(options *https.Options) {
		options.LoggerHeaderKeys = keys
	}
}

func OnStart(fn func()) https.Option {
	return func(options *https.Options) {
		options.OnStart = fn
	}
}

func OnStartFail(fn func(error)) https.Option {
	return func(options *https.Options) {
		options.OnStartFail = fn
	}
}

// OnResponded 请求响应完成后执行（请求日志已记录）
func OnResponded(fn func(ctx context.Context, data *RequestResponseData)) https.Option {
	return func(options *https.Options) {
		options.OnResponded = fn
	}
}

// OnRecovered 恢复异常时执行
func OnRecovered(fn func(ctx context.Context, data *RecoveredData)) https.Option {
	return func(options *https.Options) {
		options.OnRecovered = fn
	}
}

func OnShutdown(fn func()) https.Option {
	return func(options *https.Options) {
		options.OnShutdown = fn
	}
}

func OnExit(fn func()) https.Option {
	return func(options *https.Options) {
		options.OnExit = fn
	}
}

func Inits(inits ...https.Init) Option {
	return func(options *https.Options) {
		options.Inits = append(options.Inits, inits...)
	}
}

func Defers(defers ...https.Defer) Option {
	return func(options *https.Options) {
		options.Defers = append(options.Defers, defers...)
	}
}

func Routes(routes ...https.Route) Option {
	return func(options *https.Options) {
		options.Routes = append(options.Routes, routes...)
	}
}

// Middleware 返回可自定义优先级的中间件
// 可自定义优先级，默认优先级999(priorities仅取第一个作为中间件优先级)
func Middleware(handle https.MiddlewareHandleFunc, priorities ...int) https.Middleware {
	return https.Middleware{
		Priority: lo.FirstOr(priorities, 9999),
		Handle:   handle,
	}
}

func WrapEngineMiddleware(handle gin.HandlerFunc, priorities ...int) https.Middleware {
	return https.Middleware{
		Priority: lo.FirstOr(priorities, 9999),
		Handle: func(e *gin.Engine) gin.HandlerFunc {
			return handle
		},
	}
}

// Middlewares 注册中间件
func Middlewares(middlewares ...https.Middleware) Option {
	return func(options *https.Options) {
		options.Middlewares = append(options.Middlewares, middlewares...)
	}
}
