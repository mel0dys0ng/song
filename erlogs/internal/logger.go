package internal

import (
	"os"

	"github.com/song/metas"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	timeKey       = "ts"     // 时间字段名
	levelKey      = "lv"     // 日志级别字段名
	nameKey       = "name"   // 名称字段名
	callerKey     = "caller" // 调用者字段名
	messageKey    = "msg"    // 消息字段名
	stacktraceKey = "trace"  // 堆栈Trace字段名
)

func newZapCore(config *Config) zapcore.Core {
	hook := lumberjack.Logger{
		Filename:   config.GetFilePath(),   // 日志文件路径
		MaxSize:    config.GetMaxSize(),    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.GetMaxBuckups(), // 日志文件最多保存多少个备份
		MaxAge:     config.GetMaxAge(),     // 文件最多保存多少天
		Compress:   config.GetCompose(),    // 是否压缩
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(ToLevel(config.Level).ToZapLevel())

	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        timeKey,
		LevelKey:       levelKey,
		NameKey:        nameKey,
		CallerKey:      callerKey,
		MessageKey:     messageKey,
		StacktraceKey:  stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	var ws []zapcore.WriteSyncer
	if !metas.Mode().IsModeDebug() {
		ws = []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)}
	} else {
		ws = []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	}

	return zapcore.NewCore(
		// 编码器
		zapcore.NewJSONEncoder(encoderConfig),
		// 打印输出
		zapcore.NewMultiWriteSyncer(ws...),
		// 日志级别
		atomicLevel,
	)
}
