package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type (
	RedisStreamPublisher struct {
		*Logger

		config    redisstream.PublisherConfig
		publisher *redisstream.Publisher
	}

	RedisStreamPublisherOption func(*RedisStreamPublisher)
)

func RedisStreamPublisherLogger(opts ...LoggerOption) RedisStreamPublisherOption {
	return func(i *RedisStreamPublisher) {
		i.logger = NewLogger(opts...).logger
	}
}

func RedisStreamPublisherConfigDefaultMaxLen(n int64) RedisStreamPublisherOption {
	return func(i *RedisStreamPublisher) {
		i.config.DefaultMaxlen = n
	}
}

func RedisStreamPublisherConfigMaxLens(maxLens map[string]int64) RedisStreamPublisherOption {
	return func(i *RedisStreamPublisher) {
		i.config.Maxlens = maxLens
	}
}

func NewRedisStreamPublisher(ctx context.Context, client redis.UniversalClient, opts ...RedisStreamPublisherOption) *RedisStreamPublisher {
	if client == nil {
		erlogs.New("redis universal client is nil").Options(BaseELOptions()).PanicLog(ctx)
	}

	rs := &RedisStreamPublisher{
		config: redisstream.PublisherConfig{Client: client},
		Logger: &Logger{
			logger: watermill.NewStdLogger(false, false),
		},
	}

	for _, v := range opts {
		v(rs)
	}

	var err error
	rs.publisher, err = redisstream.NewPublisher(rs.config, rs.logger)
	if err != nil {
		erlogs.Convert(err).Wrap("failed to create publisher").Options(BaseELOptions()).PanicLog(ctx)
	}

	return rs
}

func (i *RedisStreamPublisher) Publish(ctx context.Context, topic string, msgs ...*message.Message) (err error) {
	if len(msgs) == 0 {
		err = erlogs.New("messages empty").Erorr(erlogs.OptionFields(zap.String("topic", topic)))
		return
	}

	if err = i.publisher.Publish(topic, msgs...); err != nil {
		err = erlogs.Convert(err).Wrap("failed to publish messages").Erorr(erlogs.OptionFields(zap.String("topic", topic)))
		return
	}

	erlogs.New("published messages").Options(BaseELOptions()).InfoLog(ctx,
		erlogs.OptionFields(zap.String("topic", topic), zap.Any("msgs", msgs)),
	)
	return
}

func (i *RedisStreamPublisher) GetPublisher() *redisstream.Publisher {
	return i.publisher
}

func (i *RedisStreamPublisher) Close() {
	i.publisher.Close()
}
