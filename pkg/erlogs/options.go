package erlogs

import (
	"github.com/mel0dys0ng/song/internal/core/erlogs"
	"go.uber.org/zap"
)

func Options(opts ...Option) Option {
	return erlogs.Options(opts...)
}

func OptionKindBiz() Option {
	return erlogs.OptionKind(erlogs.KindBiz)
}

func OptionKindSystem() Option {
	return erlogs.OptionKind(erlogs.KindSystem)
}

func OptionKindTrace() Option {
	return erlogs.OptionKind(erlogs.KindTrace)
}

func OptionBiz(bizID int32, bizName string) Option {
	return erlogs.OptionBiz(bizID, bizName)
}

func OptionCode(code int64) Option {
	return erlogs.OptionCode(code)
}

func OptionMsg(msg string) Option {
	return erlogs.OptionMsg(msg)
}

func OptionMsgV(args ...any) Option {
	return erlogs.OptionMsgV(args...)
}

func OptionEvent(event string) Option {
	return erlogs.OptionEvent(event)
}

func OptionContent(content string) Option {
	return erlogs.OptionContent(content)
}

func OptionContentf(format string, args ...interface{}) Option {
	return erlogs.OptionContentf(format, args...)
}

func OptionFields(fields ...zap.Field) Option {
	return erlogs.OptionFields(fields...)
}

func OptionPCs(pcs ...uintptr) Option {
	return erlogs.OptionPCs(pcs...)
}

func OptionPC(pc uintptr) Option {
	return erlogs.OptionPC(pc)
}
