package cmd

import (
	"github.com/spf13/cobra"
)

func newVersionCmd(ver string) *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "版本号",
		Long:  "显示版本号",
		Args:  cobra.MatchAll(cobra.ExactArgs(0)),
		Run: func(cmd *cobra.Command, args []string) {
			println("版本号：" + ver)
		},
	}
	return versionCmd
}
