package erlogs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mel0dys0ng/song/pkg/sys"
	"go.uber.org/zap"
)

type bufferedLog struct {
	level   Level
	message string
	fields  []zap.Field
}

var (
	logBuffer      []bufferedLog
	logBufferMutex sync.RWMutex
)

func addLogToBuffer(level Level, message string, fields []zap.Field) {
	logBufferMutex.Lock()
	defer logBufferMutex.Unlock()
	logBuffer = append(logBuffer, bufferedLog{
		level:   level,
		message: message,
		fields:  fields,
	})

	switch level {
	case LevelPanic, LevelFatal:
		sys.Panicf("%s", buildPanicMessage(message, fields))
	}
}

func buildPanicMessage(message string, fields []zap.Field) string {
	var panicMessage strings.Builder
	panicMessage.WriteString(message)
	if len(fields) > 0 {
		panicMessage.WriteString(" | fields: ")
		for i, field := range fields {
			if i > 0 {
				panicMessage.WriteString(", ")
			}
			panicMessage.WriteString(field.Key + "=")
			if field.String != "" {
				panicMessage.WriteString(field.String)
			} else {
				fmt.Fprintf(&panicMessage, "%v", field.Interface)
			}
		}
	}
	return panicMessage.String()
}

func getBufferedLogs() []bufferedLog {
	logBufferMutex.RLock()
	defer logBufferMutex.RUnlock()
	copied := make([]bufferedLog, len(logBuffer))
	copy(copied, logBuffer)
	return copied
}

func clearLogBuffer() {
	logBufferMutex.Lock()
	defer logBufferMutex.Unlock()
	logBuffer = nil
}

func logByLevel(logger *zap.Logger, level Level, message string, fields ...zap.Field) {
	switch level {
	case LevelDebug:
		logger.Debug(message, fields...)
	case LevelInfo:
		logger.Info(message, fields...)
	case LevelWarn:
		logger.Warn(message, fields...)
	case LevelError:
		logger.Error(message, fields...)
	case LevelPanic:
		logger.Panic(message, fields...)
	case LevelFatal:
		logger.Fatal(message, fields...)
	default:
		logger.Error(message, fields...)
	}
}

func flushBufferedLogs(config *Config) {
	bufferedLogs := getBufferedLogs()
	if len(bufferedLogs) > 0 {
		logger := zap.New(newZapCore(config))
		for _, log := range bufferedLogs {
			logByLevel(logger, log.level, log.message, log.fields...)
		}
		clearLogBuffer()
	}
}
