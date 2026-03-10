package cobras

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Empty struct {
	Name string
}

func NewEmptyCommand(name string) CommandInterface {
	return &Empty{Name: name}
}

// Long return the long description of command
func (c *Empty) Long() string {
	return c.Name
}

// Short return the short description of command
func (c *Empty) Short() string {
	return c.Name
}

// BindFlags bind flags
func (c *Empty) BindFlags(set *pflag.FlagSet) {

}

// The *Run functions are executed in the following order:
//   * PersistentPreRun()
//   * PreRun()
//   * Run()
//   * PostRun()
//   * PersistentPostRun()
// All functions get the same args, the arguments after the command name.

// PersistentPreRun : children of this command will inherit and execute.
func (c *Empty) PersistentPreRun(cmd *cobra.Command, args []string) {

}

// PreRun : children of this command will not inherit.
func (c *Empty) PreRun(cmd *cobra.Command, args []string) {

}

// Run : Typically the actual work function. Most commands will only implement this.
func (c *Empty) Run(cmd *cobra.Command, args []string) {

}

// PostRun : run after the Run command.
func (c *Empty) PostRun(cmd *cobra.Command, args []string) {

}

// PersistentPostRun : children of this command will inherit and execute after PostRun.
func (c *Empty) PersistentPostRun(cmd *cobra.Command, args []string) {

}
