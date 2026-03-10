package erlogs

import (
	"context"

	"go.uber.org/zap"
)

type ErLogInterface interface {
	GetLevel() Level
	GetErr() error
	GetKind() Kind
	GetBizID() int32
	GetBizName() string
	GetCode() int64
	GetMsg() string
	GetContent() string
	GetAt() int64
	GetFields() []zap.Field
	GetSkip() int
	GetPCs() []uintptr

	Options(opts []Option) *ErLog
	AppendFields(fields ...zap.Field) *ErLog
	Status(code int64, msg string) *ErLog
	Statusf(code int64, format string, args ...any) *ErLog
	UseStatusIfNot(err *ErLog) *ErLog

	Wrap(text string, fields ...zap.Field) *ErLog
	Wrapf(format string, args ...any) *ErLog
	WrapE(err error, fields ...zap.Field) *ErLog
	Unwrap() error
	Is(target error) bool
	As(target any) bool

	Debug(opts ...Option) *ErLog
	Info(opts ...Option) *ErLog
	Warn(opts ...Option) *ErLog
	Erorr(opts ...Option) *ErLog
	Panic(opts ...Option) *ErLog
	Fatal(opts ...Option) *ErLog

	DebugLog(ctx context.Context, opts ...Option)
	InfoLog(ctx context.Context, opts ...Option)
	WarnLog(ctx context.Context, opts ...Option)
	ErrorLog(ctx context.Context, opts ...Option)
	PanicLog(ctx context.Context, opts ...Option)
	FatalLog(ctx context.Context, opts ...Option)
	RecordLog(ctx context.Context, opts ...Option)

	Log(ctx context.Context, opts ...Option) *ErLog
	Clone() *ErLog
	Error() string
	String() string
}
