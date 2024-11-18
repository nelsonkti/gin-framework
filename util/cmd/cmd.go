package cmd

import "github.com/spf13/cobra"

type Script interface {
	Command() *cobra.Command
	Run(cmd *cobra.Command, args []string)
}

func Register(cmd *cobra.Command, script Script) {
	c := script.Command()
	c.Run = func(cmd *cobra.Command, args []string) {
		script.Run(cmd, args)
	}
	cmd.AddCommand(c)
}
