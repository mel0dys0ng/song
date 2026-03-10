package erlogs

import (
	"errors"

	"github.com/mel0dys0ng/song/pkg/singleton"
	"github.com/mel0dys0ng/song/pkg/sys"
)

// Initialize 初始化错误日志记录器，配置全局日志记录器
func Initialize(config *Config) {
	if config == nil {
		sys.Panic("erlogs config is nil")
	}

	singleton.Once(globalConfigKey, func() *Config {
		err := config.Check()
		if err != nil {
			sys.Panicf("failed to initialize erlogs config: %s", err.Error())
		}
		return config
	})

	flushBufferedLogs(config)
}

// New 创建一个新的 ErLog 实例，使用给定的文本作为错误消息
func New(text string, opts ...Option) ErLogInterface {
	e := Constructor(opts...)
	if text != "" {
		e.setErr(errors.New(text))
		e.setLevel(LevelError)
		e.setContent(text)
	} else {
		e.setLevel(LevelInfo)
	}
	e.capturePC()
	return e
}

// Convert 将标准 error 转换为 ErLog，如果 err 已经是 ErLog 类型则直接返回
func Convert(err error, opts ...Option) ErLogInterface {
	if err == nil {
		return nil
	}

	el, ok := errors.AsType[*ErLog](err)
	if !ok {
		e := Constructor(opts...)
		e.setErr(err)
		e.setLevel(LevelError)
		e.setContent(err.Error())
		el = e
	}

	el.capturePC()
	return el
}

// WithOptions 创建一个新的 ErLog 实例，仅应用可选参数配置
func WithOptions(opts ...Option) ErLogInterface {
	el := Constructor(opts...)
	el.capturePC()
	return el
}
