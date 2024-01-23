package cmd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
)

var runFlag = &restapp.RunFlag{}
var runFunc func(flag *restapp.RunFlag) error

func Start(config string, fun func(flag *restapp.RunFlag) error) {
	runFunc = fun
	rootCmd := newRootCmd(config)
	rootCmd.AddCommand(newStartCmd())
	rootCmd.AddCommand(newInitDbCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newStopCmd())
	rootCmd.AddCommand(newCreateSqlFileCmd())
	_ = rootCmd.Execute()
}

func newRootCmd(config string) *cobra.Command {
	name := restapp.GetAppExcName()
	var startCmd = &cobra.Command{
		Use:   name,
		Short: "启动服务",
		Long:  "启动服务",
		Run: func(cmd *cobra.Command, args []string) {
			runFlag.RunType = restapp.RunTypeStart
			if err := runFunc(runFlag); err != nil {
				panic(err)
			}
		},
	}
	startCmd.PersistentFlags().StringVar(&runFlag.Config, "config", config, "config file")
	startCmd.PersistentFlags().StringVar(&runFlag.Env, "env", "", "env 配置环境名称")
	return startCmd
}
