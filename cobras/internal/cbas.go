package internal

import (
	"context"
	"sync"

	"github.com/mel0dys0ng/song/erlogs"
	"github.com/mel0dys0ng/song/utils/sljces"
	"go.uber.org/zap"
)

var (
	_elgSys         erlogs.ErLogInterface
	_elgSysInitOnce sync.Once
)

type Cbas struct {
	commands map[string]*Command
	root     *Command
}

func New(name string) CbassInterface {
	root := NewCommand(name, nil, nil, NewEmptyCommand(name))
	return &Cbas{
		commands: map[string]*Command{root.Index: root},
		root:     root,
	}
}

func elgSys() erlogs.ErLogInterface {
	_elgSysInitOnce.Do(func() {
		_elgSys = erlogs.New(erlogs.TypeSystem(), erlogs.Log(true), erlogs.Msgf("[cobras] %s"))
	})
	return _elgSys
}

// RegisterRoot 注册根命令（替换掉默认的根命令EmptyCommand）
func (c *Cbas) RegisterRoot(newRootCommandFunc func(name string) CommandInterface) {
	root := newRootCommandFunc(c.root.Name)
	if root != nil {
		c.root = NewCommand(c.root.Name, nil, nil, root)
		c.commands[c.root.Index] = c.root
		return
	}

	elgSys().PanicL(context.Background(),
		erlogs.Msgv("failed to register root command"),
		erlogs.Content("root command is nil"),
	)
}

// RegisterCommand 为root命令添加子命令。若cmds为空，则添加Empty Command；否则添加第一个 Command。
func (c *Cbas) RegisterCommand(name string, cmds ...CommandInterface) CbaInterface {
	cmd := sljces.Index(cmds, 0, NewEmptyCommand(name))
	command := NewCommand(name, c.root, c.commands, cmd)
	c.commands[command.Index] = command
	return NewCba(c, command)
}

// RegisterExecute 注册并执行命令
func (c *Cbas) RegisterExecute(register func(c CbasInterface)) {
	if register != nil {
		register(c)
		c.Execute()
		return
	}

	elgSys().PanicL(context.Background(),
		erlogs.Msgv("failed to register and execute"),
		erlogs.Content("register func is nil"),
	)
}

// Execute execute commands
func (c *Cbas) Execute() {
	defer func() {
		if e := recover(); e != nil {
			elgSys().PanicL(context.Background(),
				erlogs.Msgv("execute commands fail"),
				erlogs.Content("execute panic"),
				erlogs.Fields(zap.Any("panicErr", e)),
			)
		}
	}()

	var root *Command
	for _, v := range c.commands {
		cmd, ok := c.commands[v.Index]
		if !ok || cmd == nil {
			elgSys().PanicL(context.Background(),
				erlogs.Msgv("command invalid"),
				erlogs.Contentf("command:%s is not exist", v.Index),
			)
			return
		}

		if len(cmd.ParentIndex) == 0 {
			root = v
			continue
		}

		pcmd, ok := c.commands[v.ParentIndex]
		if !ok || pcmd == nil || pcmd.CobraCommand == nil {
			elgSys().PanicL(context.Background(),
				erlogs.Msgv("command invalid"),
				erlogs.Contentf("command:%s is not exist", v.ParentIndex),
			)
			return
		}

		pcmd.CobraCommand.AddCommand(cmd.CobraCommand)
	}

	if root == nil || root.CobraCommand == nil {
		elgSys().PanicL(context.Background(),
			erlogs.Msgv("root command invalid"),
			erlogs.Content("root or root.CobraCommand is nil"),
		)
		return
	}

	if xer := root.CobraCommand.Execute(); xer != nil {
		elgSys().PanicL(context.Background(),
			erlogs.Msgv("failed to execute commands"),
			erlogs.Content(xer.Error()),
		)
	}
}
