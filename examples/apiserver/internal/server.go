package internal

import (
	"slices"

	"github.com/mel0dys0ng/song/core/cobras"
	"github.com/mel0dys0ng/song/core/https"
	"github.com/mel0dys0ng/song/core/metas"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type HttpServer struct {
	cobras.EmptyCommand
	metas.Options
}

func NewHttpServer(app string) cobras.CommandInterface {
	return &HttpServer{
		Options: metas.Options{
			App: app,
		},
	}
}

// Long return the long description of command
func (c *HttpServer) Long() string {
	return c.App
}

// Short return the short description of command
func (c *HttpServer) Short() string {
	return c.App
}

// BindFlags bind flags
func (c *HttpServer) BindFlags(set *pflag.FlagSet) {
	set.StringVarP(&c.Product, "product", "p", "", "the name of the product")
	set.StringVarP(&c.Config, "config", "c", "", "the path of the app config")
}

// Run : Typically the actual work function. Most commands will only implement this.
func (c *HttpServer) Run(cmd *cobra.Command, args []string) {
	// 初始化元数据
	metas.Init(&c.Options)

	// 配置项
	options := slices.Concat(
		[]https.Option{
			// 设置defer
			SetupDefers(),
			// 启动时执行
			SetupStartHook(),
			// 启动失败时执行
			SetupStartFailHook(),
			// 初始化（在配置、HTTP服务启动之前，在中间件、路由之前执行）
			SetupInits(),
			// 设置中间件
			SetupMiddlewares(),
			// 请求响应完成后执行
			SetupRespondedHook(),
			// 关闭时执行
			SetupShutdownHook(),
			// 退出时执行
			SetupExitHook(),
		},
		// 设置路由
		SetupRoutes(),
	)

	// 创建并启动HTTP服务器
	https.New(options).Serve()
}
