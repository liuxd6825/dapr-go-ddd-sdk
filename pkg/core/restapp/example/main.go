package main

import (
	"github.com/dapr/go-sdk/actor"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	restapp2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp/cmd"
)

var (
	Version   = "v1.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	restapp.GetEnvConfig()
	//config := "./config.yaml"
	cmd.StartCmd(func(flag *restapp2.RunFlag) error {
		opts := restapp2.NewRunOptions().SetFlag(flag)
		opts.AddOnStartEvent(func(server *restapp2.HttpServer) error {
			//return hserver.InitServer(server)
			return nil
		})
		runCfg := &restapp2.RunConfig{
			Subs:        subscribes,
			Controllers: controllers,
			EventTypes:  events,
			Actors:      actors,
		}
		_, err := restapp2.RunWithConfig(flag.Env, flag.Config, runCfg, opts)
		return err
	}, func(opts *cmd.Option) {
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
