package internal

import "github.com/mel0dys0ng/song/core/https"

// SetupMiddlewares 注册全局中间件
// Recovery和Logger中间件已注册，请勿重复注册
func SetupMiddlewares() https.Option {
	return https.Middlewares()
}
