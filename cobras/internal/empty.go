package internal

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type EmptyCommand struct {
	Name string
}

func NewEmptyCommand(name string) CommandInterface {
	return &EmptyCommand{Name: name}
}

// Long return the long description of command
func (c *EmptyCommand) Long() string {
	return c.Name
}

// Short return the short description of command
func (c *EmptyCommand) Short() string {
	return c.Name
}

// BindFlags bind flags
func (c *EmptyCommand) BindFlags(set *pflag.FlagSet) {

}

// The *Run functions are executed in the following order:
//   * PersistentPreRun()
//   * PreRun()
//   * Run()
//   * PostRun()
//   * PersistentPostRun()
// All functions get the same args, the arguments after the command name.

// PersistentPreRun : children of this command will inherit and execute.
func (c *EmptyCommand) PersistentPreRun(cmd *cobra.Command, args []string) {

}

// PreRun : children of this command will not inherit.
func (c *EmptyCommand) PreRun(cmd *cobra.Command, args []string) {

}

// Run : Typically the actual work function. Most commands will only implement this.
func (c *EmptyCommand) Run(cmd *cobra.Command, args []string) {

}

// PostRun : run after the Run command.
func (c *EmptyCommand) PostRun(cmd *cobra.Command, args []string) {

}

// PersistentPostRun : children of this command will inherit and execute after PostRun.
func (c *EmptyCommand) PersistentPostRun(cmd *cobra.Command, args []string) {

}
