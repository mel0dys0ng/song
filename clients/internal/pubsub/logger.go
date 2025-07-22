package pubsub

import "github.com/ThreeDotsLabs/watermill"

type (
	Logger struct {
		logger watermill.LoggerAdapter
	}

	LoggerOption struct {
		Func func(*Logger)
	}
)

func StdLogger(debug, trace bool) LoggerOption {
	return LoggerOption{
		Func: func(l *Logger) {
			l.logger = watermill.NewStdLogger(debug, trace)
		},
	}
}

func CustomLogger(logger watermill.LoggerAdapter) LoggerOption {
	return LoggerOption{
		Func: func(l *Logger) {
			l.logger = logger
		},
	}
}

func NewLogger(opts ...LoggerOption) *Logger {
	l := &Logger{
		logger: watermill.NewStdLogger(false, false),
	}

	for _, v := range opts {
		v.Func(l)
	}

	return l
}
