package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
)

var runFlag = &restapp.RunFlag{}
var runFunc func(flag *restapp.RunFlag) error

type Option struct {
	Version  string
	AppTitle string
}
type Options func(opts *Option)

func Start(config string, fun func(flag *restapp.RunFlag) error, options ...Options) {
	opts := &Option{Version: "v1.0", AppTitle: "应用服务"}
	for _, o := range options {
		if o != nil {
			o(opts)
		}
	}

	runFunc = fun
	rootCmd := newRootCmd(config, opts.AppTitle)
	rootCmd.AddCommand(newStartCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newStopCmd())
	rootCmd.AddCommand(newInitDbCmd())
	rootCmd.AddCommand(newCreateSqlFileCmd())
	rootCmd.AddCommand(newVersionCmd(opts.Version))
	rootCmd.SetVersionTemplate(opts.Version)
	rootCmd.Commands()
	_ = rootCmd.Execute()
}

func newRootCmd(config string, appTitle string) *cobra.Command {
	name := restapp.GetExeName()
	var rootCmd = &cobra.Command{
		Use:   name,
		Short: appTitle,
		Long:  appTitle,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("请使用 start, status, stop, init-db, sql-file, version, help 命令")
		},
	}
	rootCmd.PersistentFlags().StringVarP(&runFlag.Config, "config", "c", config, "配置文件名")
	rootCmd.PersistentFlags().StringVarP(&runFlag.Env, "env", "e", "", "配置文件中定义的env环境名称")
	return rootCmd
}
