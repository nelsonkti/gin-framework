package task

import (
	"fmt"
	"github.com/spf13/cobra"
)

type DemoScript struct{}

func (ds *DemoScript) Command() *cobra.Command {
	return &cobra.Command{
		Use:   "demo-script",
		Short: "修复脚本",
		Long:  ``,
	}
}

func (ds *DemoScript) Run(cmd *cobra.Command, args []string) {
	fmt.Println("执行逻辑。。。。。")
}
