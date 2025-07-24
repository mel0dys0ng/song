package internal

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug         Level = 1         // debug
	LevelInfo          Level = 2         // info
	LevelWarn          Level = 3         // warn
	LevelError         Level = 4         // error
	LevelPanic         Level = 5         // panic
	LevelFatal         Level = 6         // fatal
	LevelUnknown       Level = 7         // unknown
	LevelStringDebug         = "debug"   // debug
	LevelStringInfo          = "info"    // info
	LevelStringWarn          = "warn"    // warn
	LevelStringError         = "error"   // error
	LevelStringPanic         = "panic"   // panic
	LevelStringFatal         = "fatal"   // fatal
	LevelStringUnknown       = "unknown" // unknown
)

type Level uint8

func (lv Level) String() string {
	switch lv {
	case LevelDebug:
		return LevelStringDebug
	case LevelInfo:
		return LevelStringInfo
	case LevelWarn:
		return LevelStringWarn
	case LevelError:
		return LevelStringError
	case LevelPanic:
		return LevelStringPanic
	case LevelFatal:
		return LevelStringFatal
	default:
		return LevelStringUnknown
	}
}

func (lv Level) ToZapLevel() zapcore.Level {
	switch lv {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelPanic:
		return zapcore.PanicLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InvalidLevel
	}
}

func ToLevel(lv string) Level {
	switch strings.ToLower(lv) {
	case LevelStringDebug:
		return LevelDebug
	case LevelStringInfo:
		return LevelInfo
	case LevelStringWarn:
		return LevelWarn
	case LevelStringError:
		return LevelError
	case LevelStringPanic:
		return LevelPanic
	case LevelStringFatal:
		return LevelFatal
	default:
		return LevelUnknown
	}
}
