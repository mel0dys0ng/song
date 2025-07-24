package pubsub

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

type (
	MiddlewareRetry = middleware.Retry

	MiddlewareRetryOption struct {
		Func func(MiddlewareRetry)
	}
)

// MaxInterval 为重试的指数退避设置了上限。间隔时间不会超过 MaxInterval。
func MiddlewareRetryMaxInterval(d time.Duration) MiddlewareRetryOption {
	return MiddlewareRetryOption{
		Func: func(mr MiddlewareRetry) {
			mr.MaxInterval = d
		},
	}
}

// Multiplier 是重试之间等待间隔的缩放因子。
func MiddlewareRetryMultiplier(d float64) MiddlewareRetryOption {
	return MiddlewareRetryOption{
		Func: func(mr MiddlewareRetry) {
			mr.Multiplier = d
		},
	}
}

// MaxElapsedTime 设置了重试尝试的时间限制。如果设置为 0，则禁用此限制。
func MiddlewareRetryMaxElapsedTime(d time.Duration) MiddlewareRetryOption {
	return MiddlewareRetryOption{
		Func: func(mr MiddlewareRetry) {
			mr.MaxElapsedTime = d
		},
	}
}

// RandomizationFactor 用于在以下区间内随机化退避时间的分布：
// [当前间隔时间 * (1 - 随机化因子), 当前间隔时间 * (1 + 随机化因子)]。
func MiddlewareRetryRandomizationFactor(d float64) MiddlewareRetryOption {
	return MiddlewareRetryOption{
		Func: func(mr MiddlewareRetry) {
			mr.RandomizationFactor = d
		},
	}
}

// OnRetryHook 是一个可选函数，每次重试尝试时都会执行该函数。
// 当前重试的次数会作为 retryNum 参数传入。
func MiddlewareRetryOnRetryHook(f func(retryNum int, delay time.Duration)) MiddlewareRetryOption {
	return MiddlewareRetryOption{
		Func: func(mr MiddlewareRetry) {
			mr.OnRetryHook = f
		},
	}
}
