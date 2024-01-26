package restapp

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	"github.com/liuxd6825/dapr-go-sdk/service/common"
)

type RunConfig struct {
	AppId                  string
	HttpHost               string
	HttpPort               int
	LogLevel               applog.Level
	DaprMaxCallRecvMsgSize *int64
	DaprClient             dapr.DaprClient
	EnvConfig              *EnvConfig
}

type RegisterHandler interface {
	RegisterHandler(app *iris.Application)
}

func RunWithConfig(envName string, configFile string, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventsFunc func() []RegisterEventType, actorsFunc func() []actor.FactoryContext,
	options ...*RunOptions) (common.Service, error) {

	SetEnvName(envName)
	config, err := NewConfigByFile(configFile)
	if err != nil {
		fmt.Println(fmt.Sprintf("打开配置文件%s时出错，错误:%s", configFile, err.Error()))
		return nil, err
	}

	env := config.Env
	if len(envName) > 0 {
		env = envName
	}

	envConfig, err := config.GetEnvConfig(env)
	if err != nil {
		fmt.Println(fmt.Sprintf("获取配置环境名称[%s]时出错，错误:%s", env, err.Error()))
		return nil, err
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

	opt := NewRunOptions(options...)
	var err error

	runType := RunTypeStart
	if opt.runType != nil {
		runType = *opt.runType
	}

	switch runType {
	case RunTypeInitDB: // 是数据库初始化
		if opt.tables != nil {
			err = InitDb(opt.GetDbKey(), opt.tables, config, opt.GetPrefix())
		}
		return nil, err
	case RunTypeCreateSqlFile: // 是生成数据库脚本
		if opt.tables != nil {
			err = InitDbScript(opt.GetDbKey(), opt.tables, config, opt.GetPrefix(), opt.GetSqlFile())
		}
		return nil, err
	case RunTypeStatus: // 查看服务状态
		status(config)
		return nil, nil
	case RunTypeStop:
		_ = stop(config)
		return nil, nil
	default:
		break
	}

	//
	// 启动服务
	//

	//初始化日志
	if err = initLogs(config.Log.level, config.Log.SaveDays, config.Log.SplitHour, config.Log.LogFile, config.Log.OutputType); err != nil {
		fmt.Println(fmt.Sprintf("初始化日志文件时出错，错误:%s", err.Error()))
		return nil, err
	}

	// 启动Dapr服务
	if err = startDapr(config); err != nil {
		return nil, err
	}

	// 初始化应用
	if err = InitApplication(context.Background(), config, eventsFunc(), false, nil); err != nil {
		return nil, err
	}

	daprClient := dapr.GetDaprClient()
	runCfg := &RunConfig{
		AppId:      config.App.AppId,
		HttpHost:   config.App.HttpHost,
		HttpPort:   config.App.HttpPort,
		LogLevel:   config.Log.level,
		DaprClient: daprClient,
		EnvConfig:  config,
	}

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
