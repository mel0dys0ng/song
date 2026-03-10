package pubsub

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/mel0dys0ng/song/pkg/erlogs"
)

type (
	Messager struct {
		*Logger

		config   message.RouterConfig
		router   *message.Router
		handlers []*Handler
	}

	MessagerOption struct {
		Func func(*Messager)
	}
)

func MessagerLogger(opts ...LoggerOption) MessagerOption {
	return MessagerOption{
		Func: func(i *Messager) {
			i.logger = NewLogger(opts...).logger
		},
	}
}

func MessagerConfigCloseTimeOut(d time.Duration) MessagerOption {
	return MessagerOption{
		Func: func(i *Messager) {
			i.config = message.RouterConfig{CloseTimeout: d}
		},
	}
}

func MessagerHandlers(hs ...*Handler) MessagerOption {
	return MessagerOption{
		Func: func(i *Messager) {
			i.handlers = append(i.handlers, hs...)
		},
	}
}

func NewMessager(ctx context.Context, opts ...MessagerOption) *Messager {
	m := &Messager{
		config: message.RouterConfig{},
		Logger: &Logger{
			logger: watermill.NewStdLogger(false, false),
		},
	}

	for _, v := range opts {
		v.Func(m)
	}

	var err error
	m.router, err = message.NewRouter(m.config, m.logger)
	if err != nil {
		erlogs.Convert(err).Wrap("failed to create router").Options(BaseELOptions()).PanicLog(ctx)
	}

	return m
}

func (i *Messager) Run() {
	// 优雅关闭处理
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if len(i.handlers) == 0 {
		erlogs.New("no handler to add").Options(BaseELOptions()).PanicLog(ctx)
	}

	for _, h := range i.handlers {
		h.Validate(ctx)

		var r *message.Handler
		if len(h.PublisherTopic) == 0 || h.Publisher == nil || h.HandlerFunc == nil {
			r = i.router.AddConsumerHandler(h.Name, h.SubscriberTopic, h.Subscriber, h.NoPublishHandlerFunc)
		} else {
			r = i.router.AddHandler(h.Name, h.SubscriberTopic, h.Subscriber, h.PublisherTopic, h.Publisher, h.HandlerFunc)
		}

		if len(h.Middlewares) > 0 {
			r.AddMiddleware(h.Middlewares...)
		}
	}

	i.handlers = nil

	// 关键点：这里会阻塞运行直到收到终止信号
	if err := i.router.Run(ctx); err != nil {
		erlogs.Convert(err).Wrap("failed to run router").Options(BaseELOptions()).PanicLog(ctx)
	}
}
