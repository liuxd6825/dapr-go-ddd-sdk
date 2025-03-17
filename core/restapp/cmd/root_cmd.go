package cmd

import (
	"fmt"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/spf13/cobra"
)

var runFlag = &restapp2.RunFlag{}
var runFunc func(flag *restapp2.RunFlag) error

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

func StartCmd(fun func(flag *restapp2.RunFlag) error, options ...Options) {
	opt := &Option{Version: "1.0.0", AppTitle: "应用服务"}
	for _, o := range options {
		if o != nil {
			o(opt)
		}
	}

	restapp2.BuildTime = opt.BuildTime
	restapp2.Version = opt.Version
	restapp2.GitHead = opt.GitHead

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
	name := restapp2.GetExeName()
	rootCmd.Use = name
	rootCmd.Short = appTitle
	rootCmd.Long = appTitle
	return rootCmd
}
