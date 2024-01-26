package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
	"os"
)

func newStartCmd() *cobra.Command {
	var stopCmd = &cobra.Command{
		Use:   "start",
		Short: "启动",
		Long:  "启动应用进程与dapr守护进程",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeStart
			if err := runFunc(runFlag); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		},
	}
	return stopCmd
}
