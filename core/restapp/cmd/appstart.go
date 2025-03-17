package cmd

import (
	"github.com/dapr/go-sdk/actor"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver"
)

type AppStartOptions struct {
	Version      string                              // 版本号
	BuildTime    string                              // Build时间
	GitHead      string                              // Git地址
	AppTitle     string                              // 应用名称
	Controllers  func() []restapp2.Controller        // HTTP控制器
	Subs         func() []restapp2.RegisterSubscribe // dapr消息订阅
	Events       func() []restapp2.RegisterEventType // ddd事件订阅
	Actors       func() []actor.FactoryContext       // dapr.actor
	OnInitEvent  restapp2.OnInitEvent
	OnStartEvent restapp2.OnStartEvent
}

// StartApp
//
//	@Description: 服务启动
//	@param opts
func StartApp(opts *AppStartOptions) {
	if opts == nil {
		panic("opts is nil")
	}
	StartCmd(func(flag *restapp2.RunFlag) error {
		runOpts := restapp2.NewRunOptions().SetFlag(flag).SetTable(nil)
		runOpts.AddOnInitEvent(func(server *restapp2.HttpServer) error {
			envCfg := server.EnvConfig().App.HServer
			if envCfg.Enable {
				srcName := envCfg.SrcName
				webName := envCfg.WebName
				return hserver.InitServer(flag.MainFile, srcName, webName, server, envCfg.WatchRestart)
			}
			return nil
		})
		runOpts.AddOnStartEvent(func(server *restapp2.HttpServer) error {
			return nil
		})
		runOpts.AddOnStartEvent(opts.OnStartEvent)
		runOpts.AddOnInitEvent(opts.OnInitEvent)

		restapp2.GetSysPaths().Set("HomePath", flag.HomePath)
		_, err := restapp2.RunWithConfig(flag.Env, flag.Config, opts.Subs, opts.Controllers, opts.Events, opts.Actors, runOpts)
		return err
	}, func(o *Option) {
		o.AppTitle = opts.AppTitle
		o.Version = opts.Version
		o.BuildTime = opts.BuildTime
		o.GitHead = opts.GitHead
	})
}
