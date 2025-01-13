package main

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp/cmd"
)

var (
	Version   = "v1.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	config := "./config.yaml"
	cmd.Start(config, func(flag *restapp.RunFlag) error {
		opts := restapp.NewRunOptions().SetFlag(flag)
		opts.SetInitFunc(rs_server.InitHttpServer)
		_, err := restapp.RunWithConfig(flag.Env, flag.Config, subscribes, controllers, events, actors, opts)
		return err
	}, func(opts *cmd.Option) {
		opts.AppTitle = "XXX服务"
		opts.Version = Version
		opts.BuildTime = BuildTime
		opts.GitHead = GitHead
	})
}

// 注册消息监听器
func subscribes() []restapp.RegisterSubscribe {
	return nil
}

// 注册Http控制器
func controllers() []restapp.Controller {
	return nil
}

// 注册领域事件
func events() []restapp.RegisterEventType {
	return nil
}

// actors
// @Description: 注册actor
// @return []actor.Factory
func actors() []actor.FactoryContext {
	return nil
}
