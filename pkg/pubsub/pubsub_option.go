package pubsub

import (
	"time"

	"github.com/mel0dys0ng/song/internal/core/clients/pubsub"
)

func MessageUUID(s string) MessageOption {
	return pubsub.MessageUUID(s)
}

func MessageMetadata(d map[string]string) MessageOption {
	return pubsub.MessageMetadata(d)
}

func StdLogger(debug, trace bool) LoggerOption {
	return pubsub.StdLogger(debug, trace)
}

func CustomLogger(logger LoggerAdapter) LoggerOption {
	return pubsub.CustomLogger(logger)
}

func RedisStreamPublisherLogger(opts ...LoggerOption) RedisStreamPublisherOption {
	return pubsub.RedisStreamPublisherLogger(opts...)
}

func RedisStreamPublisherConfigDefaultMaxLen(n int64) RedisStreamPublisherOption {
	return pubsub.RedisStreamPublisherConfigDefaultMaxLen(n)
}

func RedisStreamPublisherConfigMaxLens(maxLens map[string]int64) RedisStreamPublisherOption {
	return pubsub.RedisStreamPublisherConfigMaxLens(maxLens)
}

func RedisStreamSubscriberLogger(opts ...LoggerOption) RedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberLogger(opts...)
}

func RedisStreamSubscriberConsumer(s string) RedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberConsumer(s)
}

func RedisStreamSubscriberConsumerGroup(s string) RedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberConsumerGroup(s)
}

func RedisStreamSubscriberOldestIdFirst() RedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberOldestIdFirst()
}

func RedisStreamSubscriberOldestIdLatest() RedisStreamSubscriberOption {
	return pubsub.RedisStreamSubscriberOldestIdLatest()
}

func MessagerLogger(opts ...LoggerOption) MessagerOption {
	return pubsub.MessagerLogger(opts...)
}

func MessagerConfigCloseTimeOut(d time.Duration) MessagerOption {
	return pubsub.MessagerConfigCloseTimeOut(d)
}

func MessagerHandlers(hs ...*Handler) MessagerOption {
	return pubsub.MessagerHandlers(hs...)
}

func HandlerName(s string) HandlerOption {
	return pubsub.HandlerName(s)
}

func HandlerPublisherTopic(s string) HandlerOption {
	return pubsub.HandlerPublisherTopic(s)
}

func HandlerPublisher(p Publisher) HandlerOption {
	return pubsub.HandlerPublisher(p)
}

func HandlerSubscriberTopic(s string) HandlerOption {
	return pubsub.HandlerSubscriberTopic(s)
}

func HandlerSubscriber(s Subscriber) HandlerOption {
	return pubsub.HandlerSubscriber(s)
}

func HandlerNoPublisherFunc(f NoPublishHandlerFunc) HandlerOption {
	return pubsub.HandlerNoPublisherFunc(f)
}

func HandlerMiddlewareRetry(maxRetry int, initialInterval time.Duration, opts ...MiddlewareRetryOption) HandlerOption {
	return pubsub.HandlerMiddlewareRetry(maxRetry, initialInterval, opts...)
}
