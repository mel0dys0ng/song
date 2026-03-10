package cobras

import (
	"context"

	"github.com/mel0dys0ng/song/pkg/erlogs"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type CommandRunner struct {
	commands map[string]*Command
	root     *Command
}

func New(name string) CommandRunnerInterface {
	root := NewCommand(name, NewCommandOptions{
		Cmd: NewEmptyCommand(name),
	})

	return &CommandRunner{
		root:     root,
		commands: map[string]*Command{root.Index: root},
	}
}

// RegisterRoot 注册根命令（替换掉默认的根命令EmptyCommand）
func (c *CommandRunner) RegisterRoot(newRootCommandFunc func(name string) CommandInterface) {
	root := newRootCommandFunc(c.root.Name)
	if root != nil {
		c.root = NewCommand(c.root.Name, NewCommandOptions{
			Cmd: root,
		})
		c.commands[c.root.Index] = c.root
		return
	}

	erlogs.New("failed to register root command").Options(BaseELOptions()).PanicLog(context.Background(),
		erlogs.OptionFields(zap.String("name", c.root.Name)),
	)
}

// RegisterCommand 为root命令添加子命令。若cmds为空，则添加Empty Command；否则添加第一个 Command。
func (c *CommandRunner) RegisterCommand(name string, cmds ...CommandInterface) CommandChildInterface {
	command := NewCommand(name, NewCommandOptions{
		Parent:   c.root,
		Commands: c.commands,
		Cmd:      lo.FirstOr(cmds, NewEmptyCommand(name)),
	})
	c.commands[command.Index] = command
	return NewCommandChild(c, command)
}

// RegisterExecute 注册并执行命令
func (c *CommandRunner) RegisterExecute(register func(c CommandBuilderInterface)) {
	if register != nil {
		register(c)
		c.Execute()
		return
	}

	erlogs.New("failed to register and execute").
		Options(BaseELOptions()).
		PanicLog(context.Background(),
			erlogs.OptionContent("register func is nil"),
		)
}

// Execute execute commands
func (c *CommandRunner) Execute() {
	ctx := context.Background()
	defer func() {
		if e := recover(); e != nil {
			erlogs.Newf("execute commands fail: %v", e).Options(BaseELOptions()).PanicLog(ctx)
		}
	}()

	var root *Command
	for _, v := range c.commands {
		cmd, ok := c.commands[v.Index]
		if !ok || cmd == nil {
			erlogs.New("command invalid").Options(BaseELOptions()).PanicLog(ctx,
				erlogs.OptionContentf("command:%s is not exist", v.Index),
			)
			return
		}

		if len(cmd.ParentIndex) == 0 {
			root = v
			continue
		}

		pcmd, ok := c.commands[v.ParentIndex]
		if !ok || pcmd == nil || pcmd.CobraCommand == nil {
			erlogs.New("command invalid").Options(BaseELOptions()).PanicLog(ctx,
				erlogs.OptionContentf("command:%s is not exist", v.ParentIndex),
			)
			return
		}

		pcmd.CobraCommand.AddCommand(cmd.CobraCommand)
	}

	if root == nil || root.CobraCommand == nil {
		erlogs.New("root command invalid").Options(BaseELOptions()).PanicLog(ctx,
			erlogs.OptionContent("root or root.CobraCommand is nil"),
		)
		return
	}

	if xer := root.CobraCommand.Execute(); xer != nil {
		erlogs.New("failed to execute commands").Options(BaseELOptions()).PanicLog(ctx,
			erlogs.OptionContent(xer.Error()),
		)
	}
}
