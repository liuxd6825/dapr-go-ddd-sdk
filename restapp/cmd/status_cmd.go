package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
	"os"
)

func newStatusCmd() *cobra.Command {
	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "查看状态",
		Long:  "查看应用服务状态与dapr服务的状态",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeStatus
			if err := runFunc(runFlag); err != nil {
				fmt.Println(err.Error())
				os.Exit(0)
			}
		},
	}
	return statusCmd
}
