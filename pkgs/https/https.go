package https

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/mel0dys0ng/song/pkgs/https/internal"
	"github.com/mel0dys0ng/song/pkgs/utils/cljent"
	"github.com/mel0dys0ng/song/pkgs/utils/sljces"
)

type (
	Option              = internal.Option
	PriorityMiddleware  = internal.Middleware
	Route               = internal.Route
	RequestResponseData = internal.RequestResponseData
)

const (
	ClientInfoContextKey = "X-Song-Client-Info"
)

// New return a new HTTP Server
// @Param opts []Option the option of http server
func New(opts []Option) *internal.Server {
	return internal.New(opts)
}

// GetClientInfo 获取客户端信息
func GetClientInfo(ctx *gin.Context) (res *cljent.ClientInfo) {
	v, o := ctx.Get(ClientInfoContextKey)
	if o && v != nil {
		res = v.(*cljent.ClientInfo)
	}
	return
}

// Port 设置http server端口
func Port(i int) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.Port = i
		},
	}
}

// TLSOpen 是否开启https
func TLSOpen(b bool) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.TLSOpen = b
		},
	}
}

// TLSKeyFile 设置https TLS key file path
func TLSKeyFile(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.TLSKeyFile = s
		},
	}
}

// TLSKeyCert 设置https TLS cert file path
func TLSKeyCert(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.TLSCertFile = s
		},
	}
}

func KeepAlive(b bool) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.KeepAlive = b
		},
	}
}

func ReadTimeout(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.ReadTimeout = s
		},
	}
}

func ReadHeaderTimeout(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.ReadHeaderTimeout = s
		},
	}
}

func WriteTimeout(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.WriteTimeout = s
		},
	}
}

func IdleTimeout(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.IdleTimeout = s
		},
	}
}

func HammerTime(s string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.HammerTime = s
		},
	}
}

func LoggerHeaderKeys(keys ...string) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.LoggerHeaderKeys = keys
		},
	}
}

func OnStart(fn func()) internal.Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.OnStart = fn
		},
	}
}

func OnStartFail(fn func(error)) internal.Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.OnStartFail = fn
		},
	}
}

// OnResponded 请求响应完成后执行（请求日志已记录）
func OnResponded(fn func(ctx context.Context, data *RequestResponseData)) internal.Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.OnResponded = fn
		},
	}
}

func OnShutdown(fn func()) internal.Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.OnShutdown = fn
		},
	}
}

func OnExit(fn func()) internal.Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.OnExit = fn
		},
	}
}

func Inits(inits ...internal.Init) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.Inits = append(options.Inits, inits...)
		},
	}
}

func Defers(defers ...internal.Defer) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.Defers = append(options.Defers, defers...)
		},
	}
}

func Routes(routes ...internal.Route) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.Routes = append(options.Routes, routes...)
		},
	}
}

// Middleware 返回可自定义优先级的中间件
// 可自定义优先级，默认优先级999(priorities仅取第一个作为中间件优先级)
func Middleware(handle internal.MiddlewareHandleFunc, priorities ...int) internal.Middleware {
	return internal.Middleware{
		Priority: sljces.First(priorities, 999),
		Handle:   handle,
	}
}

func WrapEngineMiddleware(handle gin.HandlerFunc, priorities ...int) internal.Middleware {
	return internal.Middleware{
		Priority: sljces.First(priorities, 999),
		Handle: func(e *gin.Engine) gin.HandlerFunc {
			return handle
		},
	}
}

// Middlewares 注册中间件
func Middlewares(middlewares ...internal.Middleware) Option {
	return internal.Option{
		Apply: func(options *internal.Options) {
			options.Middlewares = append(options.Middlewares, middlewares...)
		},
	}
}
