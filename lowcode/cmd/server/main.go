package main

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/restapp/cmd"
	"os"
)

var (
	Version   = "1.0.0"
	BuildTime = ""
	GitHead   = ""
)

var config = "./config/config.yaml"
var mainFile = "./main.html"
var workPath = ""

func main() {
	appcmd.StartCmd.Flags().StringVar(&workPath, "workPath", "", "src path")
	appcmd.StartCmd.Flags().StringVar(&config, "config", "./config/config.yaml", "config file (default is $HOME/config/config.yaml)")
	appcmd.StartCmd.Flags().StringVar(&mainFile, "mainFile", "main.html", "main file (default is $HOME/main.html)")
	appcmd.Start(config, func(flag *restapp.RunFlag) error {
		opts := restapp.NewRunOptions().SetFlag(flag).SetTable(nil)
		opts.SetInitFunc(func(server *restapp.HttpServer) error {
			appCfg := server.EnvConfig().App.RsServer
			srcName := appCfg.SrcName
			webName := appCfg.WebName
			return hserver.InitServer(mainFile, srcName, webName, server)
		})
		restapp.GetSysPaths().Set("WorkPath", getWorkPath())
		_, err := restapp.RunWithConfig(flag.Env, flag.Config, nil, nil, nil, nil, opts)
		return err
	}, func(opts *appcmd.Option) {
		opts.AppTitle = "Web服务器"
		opts.Version = Version
		opts.BuildTime = BuildTime
		opts.GitHead = GitHead
	})
}

func getWorkPath() string {
	if workPath == "" {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		workPath = path
	}
	return workPath
}
