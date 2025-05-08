package main

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/restapp/cmd"
)

var (
	Version   = "1.0.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	appcmd.Start(func(flag *restapp.RunFlag) error {
		opts := restapp.NewRunOptions().SetFlag(flag).SetTable(nil)
		opts.SetInitFunc(func(server *restapp.HttpServer) error {
			envCfg := server.EnvConfig().App.RsServer
			srcName := envCfg.SrcName
			webName := envCfg.WebName
			return hserver.InitServer(flag.MainFile, srcName, webName, server)
		})
		restapp.GetSysPaths().Set("WorkPath", flag.HomePath)
		_, err := restapp.RunWithConfig(flag.Env, flag.Config, nil, nil, nil, nil, opts)
		return err
	}, func(opts *appcmd.Option) {
		opts.AppTitle = "Web服务器"
		opts.Version = Version
		opts.BuildTime = BuildTime
		opts.GitHead = GitHead
	})
}
