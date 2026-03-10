package erlogs

import (
	"fmt"

	"go.uber.org/zap"
)

type Option func(*ErLog)

func Options(opts ...Option) Option {
	return func(e *ErLog) {
		for _, opt := range opts {
			opt(e)
		}
	}
}

func OptionKind(kind Kind) Option {
	return func(e *ErLog) {
		e.setKind(kind)
	}
}

func OptionBiz(bizID int32, bizName string) Option {
	return func(e *ErLog) {
		e.setBizID(bizID)
		e.setBizName(bizName)
	}
}

func OptionCode(code int64) Option {
	return func(e *ErLog) {
		e.setCode(code)
	}
}

func OptionMsg(msg string) Option {
	return func(e *ErLog) {
		e.setMsg(msg)
	}
}

func OptionMsgV(args ...any) Option {
	return func(e *ErLog) {
		e.setMsg(fmt.Sprintf(e.msg, args...))
	}
}

func OptionEvent(event string) Option {
	return func(e *ErLog) {
		e.appendFields(zap.String("event", event))
	}
}

func OptionContent(content string) Option {
	return func(e *ErLog) {
		e.setContent(content)
	}
}

func OptionContentf(format string, args ...interface{}) Option {
	return func(e *ErLog) {
		e.setContent(fmt.Sprintf(format, args...))
	}
}

func OptionFields(fields ...zap.Field) Option {
	return func(e *ErLog) {
		e.setFields(fields)
	}
}

// OptionSkip 设置跳过的调用栈层数
func OptionSkip(skip int) Option {
	return func(e *ErLog) {
		if skip <= 0 {
			skip = SkipDefault
		}
		e.setSkip(skip)
	}
}

func OptionPCs(pcs ...uintptr) Option {
	return func(e *ErLog) {
		e.setPCs(pcs)
	}
}

func OptionPC(pc uintptr) Option {
	return func(e *ErLog) {
		e.appendPC(pc)
	}
}
