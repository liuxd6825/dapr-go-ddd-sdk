package restapp

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	"github.com/liuxd6825/dapr-go-sdk/service/common"
)

type RunConfig struct {
	AppId                  string
	HttpHost               string
	HttpPort               int
	LogLevel               applog.Level
	DaprMaxCallRecvMsgSize *int64
	DaprClient             daprclient.DaprDddClient
	EnvConfig              *EnvConfig
}

type RegisterHandler interface {
	RegisterHandler(app *iris.Application)
}

func RunWithConfig(setEnv string, configFile string, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventsFunc func() []RegisterEventType, actorsFunc func() []actor.FactoryContext,
	options ...*RunOptions) (common.Service, error) {

	config, err := NewConfigByFile(configFile)
	if err != nil {
		panic(err)
	}

	env := config.Env
	if len(setEnv) > 0 {
		env = setEnv
	}

	envConfig, err := config.GetEnvConfig(env)
	if err != nil {
		panic(err)
	}
	return RubWithEnvConfig(envConfig, subsFunc, controllersFunc, eventsFunc, actorsFunc, options...)
}

// RubWithEnvConfig
//
//	@Description: 服务启动
//	@param config  环境配置
//	@param subsFunc  Dapr消息订阅
//	@param controllersFunc  iris服务控制器
//	@param eventsFunc DDD事件注册器
//	@param actorsFunc Actor注册器
//	@param options 启动参数， 可以根据参数启动服务，初始化数据库，生成数据库脚本等
//	@return common.Service 服务
//	@return error  错误
func RubWithEnvConfig(config *EnvConfig, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventsFunc func() []RegisterEventType, actorsFunc func() []actor.FactoryContext, options ...*RunOptions) (common.Service, error) {

	if err := InitApplication(context.Background(), config, eventsFunc(), false, nil); err != nil {
		return nil, err
	}

	opt := NewRunOptions(options...)
	// 是数据库初始化
	if opt.GetInit() {
		var err error
		if opt.tables != nil {
			err = InitDb(opt.GetDbKey(), opt.tables, config, opt.GetPrefix())
		}
		return nil, err
	}

	// 是生成数据库脚本
	if opt.GetSqlFile() != "" {
		var err error
		if opt.tables != nil {
			err = InitDbScript(opt.GetDbKey(), opt.tables, config, opt.GetPrefix(), opt.GetSqlFile())
		}
		return nil, err
	}

	daprClient := daprclient.GetDaprDDDClient()

	runCfg := &RunConfig{
		AppId:                  config.App.AppId,
		HttpHost:               config.App.HttpHost,
		HttpPort:               config.App.HttpPort,
		LogLevel:               config.Log.level,
		DaprMaxCallRecvMsgSize: config.Dapr.MaxCallRecvMsgSize,
		DaprClient:             daprClient,
		EnvConfig:              config,
	}

	fmt.Printf("---------- %s ----------\r\n", config.App.AppId)
	return run(runCfg, config.App.RootUrl, subsFunc, controllersFunc, eventsFunc, actorsFunc, options...)
}

// Run
// @Description:
// @param options
// @param app
// @param webRootPath web HttpServer URL root path
// @param subsFunc
// @param controllersFunc
// @param eventStorages
// @param eventTypesFunc
// @return error
func run(runCfg *RunConfig, webRootPath string, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventTypesFunc func() []RegisterEventType, actorsFunc func() []actor.FactoryContext,
	runOptions ...*RunOptions) (common.Service, error) {

	opt := NewRunOptions(runOptions...)

	level := runCfg.LogLevel
	if opt.level != nil {
		level = *opt.level
	}
	serverOptions := &ServiceOptions{
		AppId:      runCfg.AppId,
		HttpHost:   runCfg.HttpHost,
		HttpPort:   runCfg.HttpPort,
		LogLevel:   level,
		EventTypes: eventTypesFunc(),

		Subscribes:     subsFunc(),
		Controllers:    controllersFunc(),
		ActorFactories: actorsFunc(),
		AuthToken:      "",
		WebRootPath:    webRootPath,
		EnvConfig:      runCfg.EnvConfig,
	}

	_envConfig = runCfg.EnvConfig

	// 启动HTTP服务器
	service := NewHttpServer(runCfg.DaprClient, serverOptions)
	if err := service.Start(); err != nil {
		return service, err
	}
	return service, nil
}
