package cmd

import (
	"fmt"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/spf13/cobra"
	"os"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动",
	Long:  "启动应用进程与dapr守护进程",
	Args:  cobra.MatchAll(cobra.ExactArgs(0)),
	Run: func(cmd *cobra.Command, args []string) {
		runFlag.RunType = restapp2.RunTypeStart
		restapp2.GetSysPaths().Set("HomePath", runFlag.HomePath)
		if err := runFunc(runFlag); err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	},
}

func getStartCmd() *cobra.Command {
	homePath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	startCmd.PersistentFlags().StringVar(&runFlag.Config, "config", "./config/config.yaml", "配置文件名")
	startCmd.PersistentFlags().StringVar(&runFlag.Env, "env", "", "配置文件中定义的env环境名称")
	startCmd.PersistentFlags().StringVar(&runFlag.HomePath, "homePath", homePath, "home path")
	startCmd.PersistentFlags().StringVar(&runFlag.MainFile, "mainFile", "main.html", "main file (default is main.html)")

	return startCmd
}
