package retry

import (
	"fmt"
	"time"
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
