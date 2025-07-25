package internal

import (
	"context"
	"fmt"

	"github.com/mel0dys0ng/song/core/https"
)

// SetupStartHook http server启动时执行的函数
func SetupStartHook() https.Option {
	return https.OnStart(func() {
	})
}

// SetupStartFailHook http server 启动失败时执行的函数
func SetupStartFailHook() https.Option {
	return https.OnStartFail(func(err error) {
	})
}

// SetupShutdownHook http server关闭时执行的函数
func SetupShutdownHook() https.Option {
	return https.OnShutdown(func() {
	})
}

// SetupExitHook http server退出时执行的函数
func SetupExitHook() https.Option {
	return https.OnExit(func() {
	})
}

// SetupRespondedHook 请求响应完成后执行
func SetupRespondedHook() https.Option {
	return https.OnResponded(
		func(ctx context.Context, data *https.RequestResponseData) {
			fmt.Println("===============================")
			fmt.Printf("%+v\n", data)
			fmt.Println("===============================")
		},
	)
}
