package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/mel0dys0ng/song/pkg/erlogs"
	"go.uber.org/zap"
)

type (
	Handler struct {
		Name                 string
		PublisherTopic       string
		Publisher            message.Publisher
		SubscriberTopic      string
		Subscriber           message.Subscriber
		HandlerFunc          message.HandlerFunc
		NoPublishHandlerFunc message.NoPublishHandlerFunc
		Middlewares          []message.HandlerMiddleware
	}

	HandlerOption struct {
		Func func(*Handler)
	}
)

func HandlerName(s string) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.Name = fmt.Sprintf("%s:%s", s, watermill.NewUUID())
		},
	}
}

func HandlerPublisherTopic(s string) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.PublisherTopic = s
		},
	}
}

func HandlerPublisher(p message.Publisher) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.Publisher = p
		},
	}
}

func HandlerSubscriberTopic(s string) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.SubscriberTopic = s
		},
	}
}

func HandlerSubscriber(s message.Subscriber) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.Subscriber = s
		},
	}
}

func HandlerFunc(f message.HandlerFunc) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.HandlerFunc = f
		},
	}
}

func HandlerNoPublisherFunc(f message.NoPublishHandlerFunc) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.NoPublishHandlerFunc = f
		},
	}
}

func HandlerMiddleware(f message.HandlerMiddleware) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			i.Middlewares = append(i.Middlewares, f)
		},
	}
}

func HandlerMiddlewareRetry(maxRetry int, initialInterval time.Duration, opts ...MiddlewareRetryOption) HandlerOption {
	return HandlerOption{
		Func: func(i *Handler) {
			mr := middleware.Retry{
				MaxRetries:      maxRetry,
				InitialInterval: initialInterval,
			}

			for _, v := range opts {
				v.Func(mr)
			}

			i.Middlewares = append(i.Middlewares, mr.Middleware)
		},
	}
}

func NewHandler(ctx context.Context, opts ...HandlerOption) *Handler {
	h := &Handler{
		Middlewares: []message.HandlerMiddleware{
			middleware.Recoverer,
		},
	}

	for _, v := range opts {
		v.Func(h)
	}

	h.Validate(ctx)
	return h
}

func (h *Handler) Validate(ctx context.Context) {
	if len(h.Name) == 0 || len(h.SubscriberTopic) == 0 || h.Subscriber == nil {
		content := "option name|subscriber|subscriberTopic should not zero"
		erlogs.Unknown.PanicL(ctx, erlogs.Content(content), erlogs.Fields(zap.Any("handler", h)))
	}
}
