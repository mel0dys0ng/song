package internal

import (
	"context"
	"fmt"
	"strings"

	"github.com/mel0dys0ng/song/erlogs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const IndexDepr = ":"

type Command struct {
	// command name
	Name string
	// command index
	Index string
	// parent command index
	ParentIndex string
	// command
	Command CommandInterface
	// cobra command
	CobraCommand *cobra.Command
}

func NewCommand(name string, parent *Command, commands map[string]*Command, cmd CommandInterface) (command *Command) {
	if isCmdNil := cmd == nil; isCmdNil || len(name) == 0 {
		elgSys().PanicL(context.Background(),
			erlogs.Msgv("failed to new command"),
			erlogs.Content("arguments invalid"),
			erlogs.Fields(zap.String("name", name), zap.Bool("isCmdNil", isCmdNil)),
		)
		return
	}

	index := name
	parentIndex := ""
	var parentNames []string
	if parent != nil && len(parent.Index) > 0 {
		parentIndex = parent.Index
		parentNames = strings.Split(parentIndex, IndexDepr)
		index = fmt.Sprintf("%s%s%s", parent.Index, IndexDepr, name)
	}

	// children of this command inherit and execute.
	inheritExec := func(cb func(v *Command)) {
		for i := 1; i <= len(parentNames); i++ {
			idx := strings.Join(parentNames[:i], IndexDepr)
			if v, ok := commands[idx]; ok && v != nil && v.Command != nil {
				cb(v)
			}
		}
	}

	cc := &cobra.Command{
		Use:     fmt.Sprintf("%s [command]", name),
		Short:   cmd.Short(),
		Long:    cmd.Long(),
		PreRun:  cmd.PreRun,
		Run:     cmd.Run,
		PostRun: cmd.PostRun,
		PersistentPreRun: func(cobraCmd *cobra.Command, args []string) {
			inheritExec(func(v *Command) { v.Command.PersistentPreRun(v.CobraCommand, args) })
			cmd.PersistentPreRun(cobraCmd, args)
		},
		PersistentPostRun: func(cobraCmd *cobra.Command, args []string) {
			inheritExec(func(v *Command) { v.Command.PersistentPostRun(v.CobraCommand, args) })
			cmd.PersistentPostRun(cobraCmd, args)
		},
	}

	cmd.BindFlags(cc.Flags())

	command = &Command{
		Name:         name,
		Index:        index,
		ParentIndex:  parentIndex,
		Command:      cmd,
		CobraCommand: cc,
	}

	return
}
