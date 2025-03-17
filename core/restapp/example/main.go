package main

import (
	"github.com/dapr/go-sdk/actor"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	cmd2 "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp/cmd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp/cmd"
)

var (
	Version   = "v1.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	config := "./config.yaml"
	cmd.Start(config, func(flag *restapp2.RunFlag) error {
		opts := restapp2.NewRunOptions().SetFlag(flag)
		opts.SetInitFunc(rs_server.InitHttpServer)
		_, err := restapp2.RunWithConfig(flag.Env, flag.Config, subscribes, controllers, events, actors, opts)
		return err
	}, func(opts *cmd2.Option) {
		opts.AppTitle = "XXX服务"
		opts.Version = Version
		opts.BuildTime = BuildTime
		opts.GitHead = GitHead
	})
}

// 注册消息监听器
func subscribes() []restapp2.RegisterSubscribe {
	return nil
}

// 注册Http控制器
func controllers() []restapp2.Controller {
	return nil
}

// 注册领域事件
func events() []restapp2.RegisterEventType {
	return nil
}

// actors
// @Description: 注册actor
// @return []actor.Factory
func actors() []actor.FactoryContext {
	return nil
}
