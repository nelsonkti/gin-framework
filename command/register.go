package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-framework/command/task"
	"go-framework/util/cmd"
	"os"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "cmd",
		Short: "command line",
		Long:  ``,
	}
}

func Register() {
	command := NewCommand()

	cmd.Register(command, &task.DemoScript{})

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
