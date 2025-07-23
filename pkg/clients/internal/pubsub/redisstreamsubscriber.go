package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/redis/go-redis/v9"
)

type (
	RedisStreamSubscriber struct {
		*Logger

		config     redisstream.SubscriberConfig
		subscriber *redisstream.Subscriber
	}

	RedisStreamSubscriberOption struct {
		Func func(*RedisStreamSubscriber)
	}
)

func RedisStreamSubscriberLogger(opts ...LoggerOption) RedisStreamSubscriberOption {
	return RedisStreamSubscriberOption{
		Func: func(i *RedisStreamSubscriber) {
			i.logger = NewLogger(opts...).logger
		},
	}
}

func RedisStreamSubscriberConsumer(s string) RedisStreamSubscriberOption {
	return RedisStreamSubscriberOption{
		Func: func(i *RedisStreamSubscriber) {
			i.config.Consumer = s
		},
	}
}

func RedisStreamSubscriberConsumerGroup(s string) RedisStreamSubscriberOption {
	return RedisStreamSubscriberOption{
		Func: func(i *RedisStreamSubscriber) {
			i.config.ConsumerGroup = s
		},
	}
}

func RedisStreamSubscriberOldestIdFirst() RedisStreamSubscriberOption {
	return RedisStreamSubscriberOption{
		Func: func(i *RedisStreamSubscriber) {
			i.config.OldestId = "0"
		},
	}
}

func RedisStreamSubscriberOldestIdLatest() RedisStreamSubscriberOption {
	return RedisStreamSubscriberOption{
		Func: func(i *RedisStreamSubscriber) {
			i.config.OldestId = "$"
		},
	}
}

func NewRedisStreamSubscriber(ctx context.Context, client redis.UniversalClient, opts ...RedisStreamSubscriberOption) *RedisStreamSubscriber {
	if client == nil {
		erlogs.Unknown.PanicL(ctx, erlogs.Content("redis universal client is nil"))
	}

	rs := &RedisStreamSubscriber{
		config: redisstream.SubscriberConfig{Client: client},
		Logger: &Logger{
			logger: watermill.NewStdLogger(false, false),
		},
	}

	for _, v := range opts {
		v.Func(rs)
	}

	var err error
	rs.subscriber, err = redisstream.NewSubscriber(rs.config, rs.logger)
	if err != nil {
		erlogs.Unknown.PanicL(ctx, erlogs.ContentError(err))
	}

	return rs
}

func (i *RedisStreamSubscriber) GetSubscriber() *redisstream.Subscriber {
	return i.subscriber
}

func (i *RedisStreamSubscriber) Close() {
	i.subscriber.Close()
}
