package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mel0dys0ng/song/erlogs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	RedisStreamPublisher struct {
		*Logger

		config    redisstream.PublisherConfig
		publisher *redisstream.Publisher
	}

	RedisStreamPublisherOption struct {
		Func func(*RedisStreamPublisher)
	}
)

func RedisStreamPublisherLogger(opts ...LoggerOption) RedisStreamPublisherOption {
	return RedisStreamPublisherOption{
		Func: func(i *RedisStreamPublisher) {
			i.logger = NewLogger(opts...).logger
		},
	}
}

func RedisStreamPublisherConfigDefaultMaxLen(n int64) RedisStreamPublisherOption {
	return RedisStreamPublisherOption{
		Func: func(i *RedisStreamPublisher) {
			i.config.DefaultMaxlen = n
		},
	}
}

func RedisStreamPublisherConfigMaxLens(maxLens map[string]int64) RedisStreamPublisherOption {
	return RedisStreamPublisherOption{
		Func: func(i *RedisStreamPublisher) {
			i.config.Maxlens = maxLens
		},
	}
}

func NewRedisStreamPublisher(ctx context.Context, client redis.UniversalClient, opts ...RedisStreamPublisherOption) *RedisStreamPublisher {
	if client == nil {
		erlogs.Unknown.PanicL(ctx, erlogs.Content("redis universal client is nil"))
	}

	rs := &RedisStreamPublisher{
		config: redisstream.PublisherConfig{Client: client},
		Logger: &Logger{
			logger: watermill.NewStdLogger(false, false),
		},
	}

	for _, v := range opts {
		v.Func(rs)
	}

	var err error
	rs.publisher, err = redisstream.NewPublisher(rs.config, rs.logger)
	if err != nil {
		erlogs.Unknown.PanicL(ctx, erlogs.ContentError(err))
	}

	return rs
}

func (i *RedisStreamPublisher) Publish(ctx context.Context, topic string, msgs ...*message.Message) (err error) {
	if len(msgs) == 0 {
		err = erlogs.Unknown.Erorr(ctx, erlogs.Content("messages empty"), erlogs.Fields(zap.String("topic", topic)))
		return
	}

	if err = i.publisher.Publish(topic, msgs...); err != nil {
		err = erlogs.Unknown.Erorr(ctx, erlogs.ContentError(err), erlogs.Fields(zap.String("topic", topic)))
		return
	}

	erlogs.Ok.InfoL(ctx, erlogs.Fields(zap.String("topic", topic), zap.Any("msgs", msgs)))
	return
}

func (i *RedisStreamPublisher) GetPublisher() *redisstream.Publisher {
	return i.publisher
}

func (i *RedisStreamPublisher) Close() {
	i.publisher.Close()
}
