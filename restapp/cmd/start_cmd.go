package cmd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	var stopCmd = &cobra.Command{
		Use:   "start",
		Short: "启动",
		Long:  "启动应用进程与守护进程",
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeStart
			if err := runFunc(runFlag); err != nil {
				panic(err)
			}
		},
	}
	return stopCmd
}
