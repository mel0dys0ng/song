package cobras

type (
	CommandChildInterface interface {
		// RegisterCommand 注册命令
		RegisterCommand(name string, cmd CommandInterface)
	}

	CommandBuilderInterface interface {
		// RegisterRoot 注册根命令（替换掉默认的根命令EmptyCommand）
		RegisterRoot(newRootCommandFunc func(name string) CommandInterface)
		// RegisterCommand 注册命令，默认为EmptyCommand
		RegisterCommand(name string, cmds ...CommandInterface) CommandChildInterface
	}

	CommandRunnerInterface interface {
		CommandBuilderInterface
		// RegisterExecute 注册并执行命令
		RegisterExecute(register func(c CommandBuilderInterface))
		// Execute 执行命令
		Execute()
	}
)
