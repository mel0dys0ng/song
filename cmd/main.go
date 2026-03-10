package main

import "github.com/mel0dys0ng/song/internal/core/cobras"

func main() {
	cobras.New("song").RegisterExecute(func(c cobras.CommandBuilderInterface) {
		c.RegisterRoot(cobras.NewEmptyCommand)

		// 初始化命令
		init := c.RegisterCommand("init", nil)
		{
			// 初始化api应用
			init.RegisterCommand("api", nil)
			// 初始化job应用
			init.RegisterCommand("job", nil)
			// 初始化messaging应用
			init.RegisterCommand("messaging", nil)
			// 初始化工具
			init.RegisterCommand("tool", nil)
		}

	})
}
