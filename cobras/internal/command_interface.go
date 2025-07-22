package internal

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type CommandInterface interface {
	// Long return the long description of command
	Long() string

	// Short return the short description of command
	Short() string

	// BindFlags bind flags : children of this command will inherit and execute.
	BindFlags(set *pflag.FlagSet)

	// The *Run functions are executed in the following order:
	//   * PersistentPreRun()
	//   * PreRun()
	//   * Run()
	//   * PostRun()
	//   * PersistentPostRun()
	// All functions get the same args, the arguments after the command name.

	// PersistentPreRun : children of this command will inherit and execute.
	PersistentPreRun(cmd *cobra.Command, args []string)

	// PreRun : children of this command will not inherit.
	PreRun(cmd *cobra.Command, args []string)

	// Run : Typically the actual work function. Most commands will only implement this.
	Run(cmd *cobra.Command, args []string)

	// PostRun : run after the Run command.
	PostRun(cmd *cobra.Command, args []string)

	// PersistentPostRun : children of this command will inherit and execute after PostRun.
	PersistentPostRun(cmd *cobra.Command, args []string)
}
