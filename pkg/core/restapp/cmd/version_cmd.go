package cmd

import (
	"fmt"
	"os"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本号",
	Long:  "显示版本号",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	Run: func(cmd *cobra.Command, args []string) {
		runFlag.RunType = restapp.RunTypeVersion
		if err := runFunc(runFlag); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	},
}

func newVersionCmd(ver string) *cobra.Command {
	return VersionCmd
}
