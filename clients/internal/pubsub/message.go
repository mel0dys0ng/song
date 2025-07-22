package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type MessageOption struct {
	Func func(*message.Message)
}

func NewMessage(data string, opts ...MessageOption) *message.Message {
	m := &message.Message{
		UUID:     watermill.NewUUID(),
		Metadata: message.Metadata{},
		Payload:  message.Payload(data),
	}

	for _, v := range opts {
		v.Func(m)
	}

	return m
}

func MessageUUID(s string) MessageOption {
	return MessageOption{
		Func: func(m *message.Message) {
			m.UUID = s
		},
	}
}

func MessageMetadata(d map[string]string) MessageOption {
	return MessageOption{
		Func: func(m *message.Message) {
			m.Metadata = message.Metadata(d)
		},
	}
}
