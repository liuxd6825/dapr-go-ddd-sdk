package cmd

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/cobra"
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
	opt := &Option{Version: "1.0.0", AppTitle: "应用服务"}
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
	name := restapp.GetExeName()
	rootCmd.Use = name
	rootCmd.Short = appTitle
	rootCmd.Long = appTitle
	return rootCmd
}
