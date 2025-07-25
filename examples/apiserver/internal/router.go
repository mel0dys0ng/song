package internal

import (
	"github.com/mel0dys0ng/song/core/https"
	"github.com/mel0dys0ng/song/examples/apiserver/internal/modules/hello"
)

// SetupRoutes 路由集（模块路由集合，模块路由在模块内部定义）
func SetupRoutes() []https.Option {
	return []https.Option{
		hello.SetupRoutes("/api/hello"),
	}
}
