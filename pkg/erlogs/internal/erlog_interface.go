// Code by Song Struct2Interface. DO NOT EDIT.

package internal

import (
	"context"
)

type ErLogInterface interface {

	// Level return the level of *ErLog
	Level() string

	// Type return the type of *ErLog
	Type() string

	// Code return the code of *ErLog
	Code() int64

	// Msg return the msg of *ErLog
	Msg() string

	// Content return the content of *ErLog
	Content() string

	// At return the time of *ErLog callered
	At() int64

	// Caller return the caller of *ErLog
	Caller() string

	Format() string

	// Chain return the caller chain of *ErLog
	Chain() []*ErLog

	// OK if err == nil || err.Level <= LevelWarn return true, else return false
	OK() bool

	// WithOptions returns a new *ErLog with new options. the new option replaces the old one.
	WithOptions(opts ...Option) *ErLog

	// WithStatus returns a new *ErLog with new code and msg options.  the new option replaces the old one.
	WithStatus(code int64, msg string) *ErLog

	// WithStatusf returns a new *ErLog with new code and msg format options.  the new option replaces the old one.
	WithStatusf(code int64, format string, values ...any) *ErLog

	Debug(ctx context.Context, opts ...Option) *ErLog

	DebugL(ctx context.Context, opts ...Option)

	DebugE(ctx context.Context, opts ...Option) *ErLog

	Info(ctx context.Context, opts ...Option) *ErLog

	InfoL(ctx context.Context, opts ...Option)

	InfoE(ctx context.Context, opts ...Option) *ErLog

	Warn(ctx context.Context, opts ...Option) *ErLog

	WarnL(ctx context.Context, opts ...Option)

	WarnE(ctx context.Context, opts ...Option) *ErLog

	Erorr(ctx context.Context, opts ...Option) *ErLog

	ErorrL(ctx context.Context, opts ...Option)

	ErorrE(ctx context.Context, opts ...Option) *ErLog

	Panic(ctx context.Context, opts ...Option) *ErLog

	PanicL(ctx context.Context, opts ...Option)

	PanicE(ctx context.Context, opts ...Option) *ErLog

	Fatal(ctx context.Context, opts ...Option) *ErLog

	FatalL(ctx context.Context, opts ...Option)

	FatalE(ctx context.Context, opts ...Option) *ErLog

	RecordLog(ctx context.Context)

	Error() string

	String() string
}
