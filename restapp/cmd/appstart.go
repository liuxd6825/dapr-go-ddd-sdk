package cmd

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type AppStartOptions struct {
	Version      string                             // 版本号
	BuildTime    string                             // Build时间
	GitHead      string                             // Git地址
	AppTitle     string                             // 应用名称
	Controllers  func() []restapp.Controller        // HTTP控制器
	Subs         func() []restapp.RegisterSubscribe // dapr消息订阅
	Events       func() []restapp.RegisterEventType // ddd事件订阅
	Actors       func() []actor.FactoryContext      // dapr.actor
	OnInitEvent  restapp.OnInitEvent
	OnStartEvent restapp.OnStartEvent
}

// StartApp
//
//	@Description: 服务启动
//	@param opts
func StartApp(opts *AppStartOptions) {
	if opts == nil {
		panic("opts is nil")
	}
	StartCmd(func(flag *restapp.RunFlag) error {
		runOpts := restapp.NewRunOptions().SetFlag(flag).SetTable(nil)
		runOpts.AddOnInitEvent(func(server *restapp.HttpServer) error {
			envCfg := server.EnvConfig().App.RsServer
			if envCfg.Enable {
				srcName := envCfg.SrcName
				webName := envCfg.WebName
				return hserver.InitServer(flag.MainFile, srcName, webName, server, envCfg.WatchRestart)
			}
			return nil
		})
		runOpts.AddOnStartEvent(func(server *restapp.HttpServer) error {
			return nil
		})
		runOpts.AddOnStartEvent(opts.OnStartEvent)
		runOpts.AddOnInitEvent(opts.OnInitEvent)

		restapp.GetSysPaths().Set("HomePath", flag.HomePath)
		_, err := restapp.RunWithConfig(flag.Env, flag.Config, opts.Subs, opts.Controllers, opts.Events, opts.Actors, runOpts)
		return err
	}, func(o *Option) {
		o.AppTitle = opts.AppTitle
		o.Version = opts.Version
		o.BuildTime = opts.BuildTime
		o.GitHead = opts.GitHead
	})
}
