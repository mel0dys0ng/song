package pubsub

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mel0dys0ng/song/internal/core/clients/pubsub"
	"github.com/redis/go-redis/v9"
)

type (
	Logger                      = pubsub.Logger
	LoggerOption                = pubsub.LoggerOption
	LoggerAdapter               = watermill.LoggerAdapter
	MessageOption               = pubsub.MessageOption
	RedisStreamPublisher        = pubsub.RedisStreamPublisher
	RedisStreamPublisherOption  = pubsub.RedisStreamPublisherOption
	RedisStreamSubscriber       = pubsub.RedisStreamSubscriber
	RedisStreamSubscriberOption = pubsub.RedisStreamSubscriberOption
	Messager                    = pubsub.Messager
	MessagerOption              = pubsub.MessagerOption
	Publisher                   = message.Publisher
	Subscriber                  = message.Subscriber
	Handler                     = pubsub.Handler
	HandlerOption               = pubsub.HandlerOption
	Message                     = message.Message
	NoPublishHandlerFunc        = message.NoPublishHandlerFunc
	MiddlewareRetryOption       = pubsub.MiddlewareRetryOption
)

func NewMessage(data string, opts ...MessageOption) *Message {
	return pubsub.NewMessage(data, opts...)
}

func NewMessager(ctx context.Context, opts ...MessagerOption) *Messager {
	return pubsub.NewMessager(ctx, opts...)
}

func NewHandler(ctx context.Context, opts ...HandlerOption) *Handler {
	return pubsub.NewHandler(ctx, opts...)
}

func NewRedisStreamPublisher(ctx context.Context, client redis.UniversalClient, opts ...RedisStreamPublisherOption) *RedisStreamPublisher {
	return pubsub.NewRedisStreamPublisher(ctx, client, opts...)
}

func NewRedisStreamSubscriber(ctx context.Context, client redis.UniversalClient, opts ...RedisStreamSubscriberOption) *RedisStreamSubscriber {
	return pubsub.NewRedisStreamSubscriber(ctx, client, opts...)
}

func MiddlewareRetryMaxInterval(d time.Duration) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryMaxInterval(d)
}

func MiddlewareRetryMultiplier(d float64) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryMultiplier(d)
}

func MiddlewareRetryMaxElapsedTime(d time.Duration) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryMaxElapsedTime(d)
}

func MiddlewareRetryRandomizationFactor(d float64) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryRandomizationFactor(d)
}

func MiddlewareRetryOnRetryHook(f func(retryNum int, delay time.Duration)) pubsub.MiddlewareRetryOption {
	return pubsub.MiddlewareRetryOnRetryHook(f)
}
