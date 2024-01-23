package cmd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
)

func newStatusCmd() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "查看状态",
		Long:  "查看应用服务状态与dapr服务的状态",
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeStatus
			if err := runFunc(runFlag); err != nil {
				panic(err)
			}
		},
	}
	return statusCmd
}
