package restapp

import (
	"fmt"
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/service/common"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type StartOptions struct {
	AppId      string
	HttpHost   string
	HttpPort   int
	LogLevel   applog.Level
	DaprClient daprclient.DaprDddClient
}

type RegisterSubscribe interface {
	GetSubscribes() *[]ddd.Subscribe
	GetHandler() ddd.QueryEventHandler
}

type registerSubscribe struct {
	subscribes *[]ddd.Subscribe
	handler    ddd.QueryEventHandler
}

func NewRegisterSubscribe(subscribes *[]ddd.Subscribe, handler ddd.QueryEventHandler) RegisterSubscribe {
	return &registerSubscribe{
		subscribes: subscribes,
		handler:    handler,
	}
}

func (r *registerSubscribe) GetSubscribes() *[]ddd.Subscribe {
	return r.subscribes
}

func (r *registerSubscribe) GetHandler() ddd.QueryEventHandler {
	return r.handler
}

type RegisterController struct {
	RelativePath string
	Controllers  []interface{}
}

type Controller interface {
	BeforeActivation(b mvc.BeforeActivation)
}

type RegisterEventType struct {
	EventType string
	Revision  string
	NewFunc   ddd.NewEventFunc
}

var EmptyActors = func() *[]actor.Factory {
	return &[]actor.Factory{}
}

var DddActors = func() *[]actor.Factory {
	return &[]actor.Factory{
		aggregateSnapshotActorFactory,
	}
}

func aggregateSnapshotActorFactory() actor.Server {
	client, err := daprclient.GetDaprDDDClient().DaprClient()
	if err != nil {
		panic(err)
	}
	return ddd.NewAggregateSnapshotActorService(client)
}

func RunWithConfig(envType string, configFile string, subsFunc func() *[]RegisterSubscribe,
	controllersFunc func() *[]Controller, eventsFunc func() *[]RegisterEventType, actorsFunc func() *[]actor.Factory) (common.Service, error) {
	config, err := NewConfigByFile(configFile)
	if err != nil {
		panic(err)
	}

	envTypeValue := config.EnvType
	if len(envType) > 0 {
		envTypeValue = envType
	}

	envConfig, err := config.GetEnvConfig(envTypeValue)
	if err != nil {
		panic(err)
	}
	return RubWithEnvConfig(envConfig, subsFunc, controllersFunc, eventsFunc, actorsFunc)
}

func RubWithEnvConfig(config *EnvConfig, subsFunc func() *[]RegisterSubscribe,
	controllersFunc func() *[]Controller, eventsFunc func() *[]RegisterEventType, actorsFunc func() *[]actor.Factory) (common.Service, error) {
	if !config.Mongo.IsEmpty() {
		initMongo(&config.Mongo)
	}

	//创建dapr客户端
	daprClient, err := daprclient.NewDaprDddClient(config.Dapr.GetHost(), config.Dapr.GetHttpPort(), config.Dapr.GetGrpcPort())
	if err != nil {
		panic(err)
	}

	daprclient.SetDaprDddClient(daprClient)

	options := &StartOptions{
		AppId:      config.App.AppId,
		HttpHost:   config.App.HttpHost,
		HttpPort:   config.App.HttpPort,
		LogLevel:   config.Log.GetLevel(),
		DaprClient: daprClient,
	}

	//创建dapr事件存储器
	esMap := make(map[string]ddd.EventStorage)
	for _, pubsubName := range config.Dapr.Pubsubs {
		eventStorage, err := ddd.NewGrpcEventStorage(daprClient, ddd.PubsubName(pubsubName))
		if err != nil {
			panic(err)
		}
		esMap[pubsubName] = eventStorage
		esMap[""] = eventStorage
	}

	return Run(options, config.App.RootUrl, subsFunc, controllersFunc, esMap, eventsFunc, actorsFunc)
}

//
// Run
// @Description:
// @param options
// @param app
// @param webRootPath web service URL root path
// @param subsFunc
// @param controllersFunc
// @param eventStorages
// @param eventTypesFunc
// @return error
//
func Run(options *StartOptions, webRootPath string, subsFunc func() *[]RegisterSubscribe,
	controllersFunc func() *[]Controller, eventStorages map[string]ddd.EventStorage,
	eventTypesFunc func() *[]RegisterEventType, actorsFunc func() *[]actor.Factory) (common.Service, error) {

	fmt.Printf("---------- %s ----------\r\n", options.AppId)
	ddd.Init(options.AppId)
	applog.Init(options.DaprClient, options.AppId, options.LogLevel)

	serverOptions := &ServiceOptions{
		AppId:          options.AppId,
		HttpHost:       options.HttpHost,
		HttpPort:       options.HttpPort,
		LogLevel:       options.LogLevel,
		EventTypes:     eventTypesFunc(),
		EventStorages:  eventStorages,
		Subscribes:     subsFunc(),
		Controllers:    controllersFunc(),
		ActorFactories: actorsFunc(),
		AuthToken:      "",
		WebRootPath:    webRootPath,
	}
	service := NewService(options.DaprClient, serverOptions)
	if err := service.Start(); err != nil {
		return service, err
	}
	return service, nil
}
