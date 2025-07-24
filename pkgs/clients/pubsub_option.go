package clients

import (
	"time"

	"github.com/mel0dys0ng/song/pkgs/clients/internal/pubsub"
)

func PubSubMessageUUID(s string) PubSubMessageOpption {
	return pubsub.MessageUUID(s)
}

func PubSubMessageMetadata(d map[string]string) PubSubMessageOpption {
	return pubsub.MessageMetadata(d)
}

func PubSubStdLogger(debug, trace bool) PubSubLoggerOption {
	return pubsub.StdLogger(debug, trace)
}

func PubSubCustomLogger(logger PubSubLoggerAdapter) PubSubLoggerOption {
	return pubsub.CustomLogger(logger)
}

func PubSubRedisStreamPublisherLogger(opts ...PubSubLoggerOption) PubSubRedisStreamPublisherOption {
	return pubsub.RedisStreamPublisherLogger(opts...)
}

func PubSubRedisStreamPublisherConfigDefaultMaxLen(n int64) PubSubRedisStreamPublisherOption {
	return pubsub.RedisStreamPublisherConfigDefaultMaxLen(n)
}

func PubSubRedisStreamPublisherConfigMaxLens(maxLens map[string]int64) PubSubRedisStreamPublisherOption {
	return pubsub.RedisStreamPublisherConfigMaxLens(maxLens)
}

func PubSubRedisStreamSubscriberLogger(opts ...PubSubLoggerOption) PubSubRedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberLogger(opts...)
}

func PubSubRedisStreamSubscriberConsumer(s string) PubSubRedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberConsumer(s)
}

func PubSubRedisStreamSubscriberConsumerGroup(s string) PubSubRedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberConsumerGroup(s)
}

func PubSubRedisStreamSubscriberOldestIdFirst() PubSubRedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberOldestIdFirst()
}

func PubSubRedisStreamSubscriberOldestIdLatest() PubSubRedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberOldestIdLatest()
}

func PubSubMessagerLogger(opts ...PubSubLoggerOption) PubSubMessagerOption {
	return pubsub.MessagerLogger(opts...)
}

func PubSubMessagerConfigCloseTimeOut(d time.Duration) PubSubMessagerOption {
	return pubsub.MessagerConfigCloseTimeOut(d)
}

func PubSubMessagerHandlers(hs ...*PubSubHandler) PubSubMessagerOption {
	return pubsub.MessagerHandlers(hs...)
}

func PubSubHandlerName(s string) PubSubHandlerOption {
	return pubsub.HandlerName(s)
}

func PubSubHandlerPublisherTopic(s string) PubSubHandlerOption {
	return pubsub.HandlerPublisherTopic(s)
}

func PubSubHandlerPublisher(p PubSubPublisher) PubSubHandlerOption {
	return pubsub.HandlerPublisher(p)
}

func PubSubHandlerSubscriberTopic(s string) PubSubHandlerOption {
	return pubsub.HandlerSubscriberTopic(s)
}

func PubSubHandlerSubscriber(s PubSubSubscriber) PubSubHandlerOption {
	return pubsub.HandlerSubscriber(s)
}

func PubSubHandlerHasPublisherFunc(f PubSubHandlerFunc) PubSubHandlerOption {
	return pubsub.HandlerFunc(f)
}

func PubSubHandlerNoPublisherFunc(f PubSubNoPublishHandlerFunc) PubSubHandlerOption {
	return pubsub.HandlerNoPublisherFunc(f)
}

func PubSubHandlerMiddleware(m PubSubMiddleware) PubSubHandlerOption {
	return pubsub.HandlerMiddleware(m)
}

func PubSubHandlerMiddlewareRetry(maxRetry int, initialInterval time.Duration, opts ...pubsub.MiddlewareRetryOption) PubSubHandlerOption {
	return pubsub.HandlerMiddlewareRetry(maxRetry, initialInterval, opts...)
}

// MaxInterval 为重试的指数退避设置了上限。间隔时间不会超过 MaxInterval。
func PubSubMiddlewareRetryMaxInterval(d time.Duration) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryMaxInterval(d)
}

// Multiplier 是重试之间等待间隔的缩放因子。
func PubSubMiddlewareRetryMultiplier(d float64) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryMultiplier(d)
}

// MaxElapsedTime 设置了重试尝试的时间限制。如果设置为 0，则禁用此限制。
func PubSubMiddlewareRetryMaxElapsedTime(d time.Duration) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryMaxElapsedTime(d)
}

// RandomizationFactor 用于在以下区间内随机化退避时间的分布：
// [当前间隔时间 * (1 - 随机化因子), 当前间隔时间 * (1 + 随机化因子)]。
func PubSubMiddlewareRetryRandomizationFactor(d float64) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryRandomizationFactor(d)
}

// OnRetryHook 是一个可选函数，每次重试尝试时都会执行该函数。
// 当前重试的次数会作为 retryNum 参数传入。
func PubSubMiddlewareRetryOnRetryHook(f func(retryNum int, delay time.Duration)) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryOnRetryHook(f)
}
