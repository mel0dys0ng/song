package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/mel0dys0ng/song/utils/result"
	"golang.org/x/sync/singleflight"
)

const (
	NumDefault   = 3                     // 重试次数
	DelayDefault = 20 * time.Millisecond // 重试间隔时间
)

var (
	sg = &singleflight.Group{}
)

type (
	retry struct {
		num             uint32        // 重试最大次数
		delay           time.Duration // 重试间隔时间
		singleflightKey string        // singleflight Key，默认为空，不启用；若不为空，则使用该Key启用singleflight
	}

	Option struct {
		apply func(r *retry)
	}

	HandlerFunc[T any] func(ctx context.Context) result.Interface[T]
)

// SingleflightKey 设置singleflight Key，默认为空，不启用；若不为空，则使用该Key启用singleflight
func SingleflightKey(s string) Option {
	return Option{
		apply: func(r *retry) {
			r.singleflightKey = s
		},
	}
}

// SingleflightKeyf 设置singleflight Key，默认为空，不启用；若不为空，则使用该Key启用singleflight
func SingleflightKeyf(format string, values ...any) Option {
	return Option{
		apply: func(r *retry) {
			r.singleflightKey = fmt.Sprintf(format, values...)
		},
	}
}

// Delay 设置延迟重试时间
func Delay(d time.Duration) Option {
	return Option{
		apply: func(r *retry) {
			r.delay = d
		},
	}
}

// DelaySecond 设置延迟重试时间，秒级时间
func DelaySecond(d int64) Option {
	return Option{
		apply: func(r *retry) {
			r.delay = time.Duration(d) * time.Second
		},
	}
}

// DelayMillisecond 设置延迟重试时间，毫秒级时间
func DelayMillisecond(d int64) Option {
	return Option{
		apply: func(r *retry) {
			r.delay = time.Duration(d) * time.Millisecond
		},
	}
}

// Num 设置重试次数
func Num(i uint32) Option {
	return Option{
		apply: func(r *retry) {
			r.num = i
		},
	}
}

// Do 执行handler，若执行失败后进行重试。
// 重试间隔按照delay时间的1/2递增（delay>0时）。
// 默认值: Num=3, Delay=10ms, SingleflightKey=""。
// res = handler(ctx)。
func Do[T any](ctx context.Context, handler HandlerFunc[T], opts ...Option) (res result.Interface[T]) {
	rt := newRetry(opts)

	if len(rt.singleflightKey) > 0 {
		r, _, _ := sg.Do(rt.singleflightKey,
			func() (res any, err error) { res = do(ctx, rt, handler); return },
		)
		return r.(result.Interface[T])
	}

	return do(ctx, rt, handler)
}

func newRetry(opts []Option) *retry {
	rt := &retry{
		num:             NumDefault,
		delay:           DelayDefault,
		singleflightKey: "",
	}

	for _, v := range opts {
		if v.apply != nil {
			v.apply(rt)
		}
	}

	if rt.num <= 0 {
		rt.num = NumDefault
	}

	if rt.delay <= 0 {
		rt.delay = DelayDefault
	}

	return rt
}

func do[T any](ctx context.Context, rt *retry, handler HandlerFunc[T]) (res result.Interface[T]) {
	// 检查 rt 是否为 nil
	if rt == nil {
		return result.Error[T](fmt.Errorf("retry config is nil"))
	}

	// 捕获 panic
	defer func() {
		if er := recover(); er != nil {
			res = result.Error[T](fmt.Errorf("panic occurred: %v", er))
		}
	}()

	for i := uint32(0); i < rt.num; i++ {
		if res = handler(ctx); res.Err() == nil {
			break
		}

		if rt.delay > 0 {
			time.Sleep(rt.delay * time.Duration((i+2)/2))
		}
	}

	return
}
