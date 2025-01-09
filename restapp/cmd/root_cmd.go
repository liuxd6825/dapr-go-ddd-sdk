package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
	"os"
)

var runFlag = &restapp.RunFlag{}
var runFunc func(flag *restapp.RunFlag) error

type Option struct {
	Version   string
	AppTitle  string
	BuildTime string
	GitHead   string
}
type Options func(opts *Option)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("请使用 start, status, stop, init-db, sql-file, version, help 命令")
	},
}

func Start(fun func(flag *restapp.RunFlag) error, options ...Options) {
	opt := &Option{Version: "", AppTitle: "应用服务"}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}

	restapp.BuildTime = opt.BuildTime
	restapp.Version = opt.Version
	restapp.GitHead = opt.GitHead

	runFunc = fun
	rootCmd = newRootCmd(opt.AppTitle)
	rootCmd.AddCommand(getStartCmd())
	rootCmd.AddCommand(newStatusCmd())
	rootCmd.AddCommand(newStopCmd())
	rootCmd.AddCommand(newInitDbCmd())
	rootCmd.AddCommand(newCreateSqlFileCmd())
	rootCmd.AddCommand(newVersionCmd(opt.Version))
	rootCmd.SetVersionTemplate(opt.Version)
	rootCmd.Commands()
	_ = rootCmd.Execute()
}

func newRootCmd(appTitle string) *cobra.Command {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	name := restapp.GetExeName()
	rootCmd.Use = name
	rootCmd.Short = appTitle
	rootCmd.Long = appTitle
	rootCmd.PersistentFlags().StringVarP(&runFlag.Config, "config", "c", "./config/config.yaml", "配置文件名")
	rootCmd.PersistentFlags().StringVarP(&runFlag.Env, "env", "e", "", "配置文件中定义的env环境名称")
	rootCmd.PersistentFlags().StringVarP(&runFlag.WorkPath, "workPath", "w", workPath, "src path")
	rootCmd.PersistentFlags().StringVarP(&runFlag.MainFile, "mainFile", "m", "main.html", "main file (default is $HOME/main.html)")
	return rootCmd
}
