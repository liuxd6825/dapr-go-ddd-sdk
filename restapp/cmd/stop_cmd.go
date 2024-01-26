package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
	"os"
)

func newStopCmd() *cobra.Command {
	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "停止服务",
		Long:  "停止应用进程与守护进程",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeStop
			if err := runFunc(runFlag); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		},
	}
	return stopCmd
}
