package clients

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mel0dys0ng/song/pkg/clients/internal/pubsub"
	"github.com/redis/go-redis/v9"
)

type (
	PubSubLogger                      = pubsub.Logger
	PubSubLoggerOption                = pubsub.LoggerOption
	PubSubLoggerAdapter               = watermill.LoggerAdapter
	PubSubMessageOpption              = pubsub.MessageOption
	PubSubRedisStreamPublisher        = pubsub.RedisStreamPublisher
	PubSubRedisStreamPublisherOption  = pubsub.RedisStreamPublisherOption
	PubSubRedisStreamSubscriber       = pubsub.RedisStreamSubscriber
	PubSubRedisStreamSubscriberOption = pubsub.RedisStreamSubscriberOption
	PubSubMessager                    = pubsub.Messager
	PubSubMessagerOption              = pubsub.MessagerOption
	PubSubPublisher                   = message.Publisher
	PubSubSubscriber                  = message.Subscriber
	PubSubHandler                     = pubsub.Handler
	PubSubHandlerOption               = pubsub.HandlerOption
	PubSubMessage                     = message.Message
	PubSubNoPublishHandlerFunc        = message.NoPublishHandlerFunc
	PubSubHandlerFunc                 = message.HandlerFunc
	PubSubMiddleware                  = message.HandlerMiddleware
)

func NewPubSubMessage(data string, opts ...PubSubMessageOpption) *PubSubMessage {
	return pubsub.NewMessage(data, opts...)
}

func NewPubSubMessager(ctx context.Context, opts ...PubSubMessagerOption) *PubSubMessager {
	return pubsub.NewMessager(ctx, opts...)
}

func NewPubSubHandler(ctx context.Context, opts ...PubSubHandlerOption) *PubSubHandler {
	return pubsub.NewHandler(ctx, opts...)
}

func NewPubSubRedisStreamPublisher(ctx context.Context, client redis.UniversalClient, opts ...PubSubRedisStreamPublisherOption) *PubSubRedisStreamPublisher {
	return pubsub.NewRedisStreamPublisher(ctx, client)
}

func NewPubSubRedisStreamSubscriber(ctx context.Context, client redis.UniversalClient, opts ...PubSubRedisStreamSubscriber) *PubSubRedisStreamSubscriber {
	return pubsub.NewRedisStreamSubscriber(ctx, client)
}
