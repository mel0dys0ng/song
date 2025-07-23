package erlogs

import (
	"context"
	"errors"

	"github.com/mel0dys0ng/song/pkg/erlogs/internal"
)

type (
	ErLogInterface = internal.ErLogInterface
	Config         = internal.Config
)

func New(opts ...internal.Option) ErLogInterface {
	return internal.New(opts)
}

func With(err error, opts ...internal.Option) ErLogInterface {
	return Convert(err).WithOptions(opts...)
}

func Convert(err error) ErLogInterface {
	if err == nil {
		return Ok
	}

	var v *internal.ErLog
	if errors.As(err, &v) {
		return v
	}

	return Unknown.WithOptions(ContentError(err))
}

// OK if err == nil || err.Level <= LevelWarn return true, else return false
func OK(err error) bool {
	return err == nil || Convert(err).OK()
}

// StartTrace 开启一个Trace Span，开始追踪.
// 和EndTarce配合追踪一个函数或方法.
// 适用于可自定义的函数或方法内部，放在函数或方法的第一行.
// 示例：
//
//	func Test(ctx context.Context) (err error) {
//			ctx = erlogs.StartTrace(ctx, "name")
//			defer erlogs.EndTrace(ctx, err)
//			// some codes...
//			return
//	}
func StartTrace(ctx context.Context, name string) context.Context {
	return internal.StartTrace(ctx, name)
}

// EndTrace 结束一个Trace Span，结束追踪.
// 和StartTrace配合追踪一个函数或方法.
// 适用于可自定义的函数或方法内部，放在函数或方法除了StartTrace外的第一行第一个defer里
// 示例：
//
//	func Test(ctx context.Context) (err error) {
//			ctx = erlogs.StartTrace(ctx, "name")
//			defer erlogs.EndTrace(ctx, err)
//			// some codes...
//			return
//	}
func EndTrace(ctx context.Context, err error) {
	span := internal.EndTrace(ctx)
	if span != nil {
		if err != nil {
			With(err, Log(true), TypeTrace()).InfoL(ctx, Fields(span.ZapFields()...))
		} else {
			Ok.InfoL(ctx, TypeTrace(), Fields(span.ZapFields()...))
		}
	}
}

// Trace 追踪回调函数fn.
// 适用于无法定制的函数或方法，比如第三方的接口请求、mysql或redis的代码执行.
// 示例：
//
//	req := &TestRequest{}
//	err := erlogs.Trace(ctx, "traceName", func(ctx context.Context) (err error) {
//		// some codes ...
//		err = Test(ctx, req)
//		return
//	})
func Trace(ctx context.Context, name string, fn func(ctx context.Context) error) (err error) {
	ctx = internal.StartTrace(ctx, name)
	defer EndTrace(ctx, err)
	return fn(ctx)
}

// TraceSpanFromContext 从context中获取trace span
func TraceSpanFromContext(ctx context.Context) *internal.TraceSpan {
	return internal.TraceSpanFromContext(ctx)
}
