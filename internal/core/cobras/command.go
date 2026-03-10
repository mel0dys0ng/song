package cobras

import (
	"context"
	"fmt"
	"strings"

	"github.com/mel0dys0ng/song/pkg/erlogs"
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

type NewCommandOptions struct {
	// parent command
	Parent *Command
	// command map
	Commands map[string]*Command
	// command
	Cmd CommandInterface
}

// NewCommand new command
func NewCommand(name string, opts NewCommandOptions) (command *Command) {
	if isCmdNil := opts.Cmd == nil; isCmdNil || len(name) == 0 {
		erlogs.New("failed to new command").Options(BaseELOptions()).PanicLog(context.Background(),
			erlogs.OptionFields(
				zap.String("name", name),
				zap.Bool("isCmdNil", isCmdNil),
				zap.Any("opts", opts),
			),
		)
		return
	}

	index := name
	parentIndex := ""
	var parentNames []string
	if opts.Parent != nil && len(opts.Parent.Index) > 0 {
		parentIndex = opts.Parent.Index
		parentNames = strings.Split(parentIndex, IndexDepr)
		index = fmt.Sprintf("%s%s%s", opts.Parent.Index, IndexDepr, name)
	}

	// children of this command inherit and execute.
	inheritExec := func(cb func(v *Command)) {
		for i := 1; i <= len(parentNames); i++ {
			idx := strings.Join(parentNames[:i], IndexDepr)
			if v, ok := opts.Commands[idx]; ok && v != nil && v.Command != nil {
				cb(v)
			}
		}
	}

	cc := &cobra.Command{
		Use:     fmt.Sprintf("%s [command]", name),
		Short:   opts.Cmd.Short(),
		Long:    opts.Cmd.Long(),
		PreRun:  opts.Cmd.PreRun,
		Run:     opts.Cmd.Run,
		PostRun: opts.Cmd.PostRun,
		PersistentPreRun: func(cobraCmd *cobra.Command, args []string) {
			inheritExec(func(v *Command) { v.Command.PersistentPreRun(v.CobraCommand, args) })
			opts.Cmd.PersistentPreRun(cobraCmd, args)
		},
		PersistentPostRun: func(cobraCmd *cobra.Command, args []string) {
			inheritExec(func(v *Command) { v.Command.PersistentPostRun(v.CobraCommand, args) })
			opts.Cmd.PersistentPostRun(cobraCmd, args)
		},
	}

	opts.Cmd.BindFlags(cc.Flags())

	command = &Command{
		Name:         name,
		Index:        index,
		ParentIndex:  parentIndex,
		Command:      opts.Cmd,
		CobraCommand: cc,
	}

	return
}
