package internal

import (
	"github.com/mel0dys0ng/song/core/cobras"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Root struct {
	cobras.EmptyCommand
}

func NewRoot(app string) cobras.CommandInterface {
	return &Root{}
}

// Long return the long description of command
func (c *Root) Long() string {
	return ""
}

// Short return the short description of command
func (c *Root) Short() string {
	return ""
}

// BindFlags bind flags
func (c *Root) BindFlags(set *pflag.FlagSet) {
}

// Run : Typically the actual work function. Most commands will only implement this.
func (c *Root) Run(cmd *cobra.Command, args []string) {

}
