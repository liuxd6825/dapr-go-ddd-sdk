package restapp

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/go-sdk/actor"
	"github.com/liuxd6825/go-sdk/service/common"
)

type RunConfig struct {
	AppId      string
	HttpHost   string
	HttpPort   int
	LogLevel   applog.Level
	DaprClient daprclient.DaprDddClient
}

type RunOptions struct {
	tables *Tables
	initDb *bool
	prefix *string
	level  *logs.Level
}

type RegisterSubscribe interface {
	GetSubscribes() []*ddd.Subscribe
	GetHandler() ddd.QueryEventHandler
	GetInterceptor() []ddd.SubscribeInterceptorFunc
}

type registerSubscribe struct {
	subscribes   []*ddd.Subscribe
	handler      ddd.QueryEventHandler
	interceptors []ddd.SubscribeInterceptorFunc
}

type RegisterController struct {
	RelativePath string
	Controllers  []interface{}
}

type Controller interface {
	BeforeActivation(b mvc.BeforeActivation)
}

type RegisterHandler interface {
	RegisterHandler(app *iris.Application)
}

type RegisterEventType struct {
	EventType string
	Version   string
	NewFunc   ddd.NewEventFunc
}

var _actorsFactory []actor.Factory = []actor.Factory{
	aggregateSnapshotActorFactory,
}

func NewRunOptions(opts ...*RunOptions) *RunOptions {
	o := &RunOptions{
		tables: nil,
	}
	for _, item := range opts {
		if item.tables != nil {
			o.tables = item.tables
		}
		if item.initDb != nil {
			o.initDb = item.initDb
		}
		if item.prefix != nil {
			o.prefix = item.prefix
		}
	}
	return o
}

type RegisterSubscribeOptions struct {
	interceptors []ddd.SubscribeInterceptorFunc
}

func (o *RegisterSubscribeOptions) SetInterceptors(v []ddd.SubscribeInterceptorFunc) *RegisterSubscribeOptions {
	o.interceptors = v
	return o
}

var _subscribeInterceptor []ddd.SubscribeInterceptorFunc

func RegisterSubscribeInterceptor(items ...ddd.SubscribeInterceptorFunc) {
	_subscribeInterceptor = append(_subscribeInterceptor, items...)
}

func NewRegisterSubscribeOptions(opts ...*RegisterSubscribeOptions) *RegisterSubscribeOptions {
	o := &RegisterSubscribeOptions{}
	for _, item := range opts {
		if item.interceptors != nil {
			o.interceptors = item.interceptors
		}
	}
	return o
}

func NewRegisterSubscribe(subscribes *[]ddd.Subscribe, handler ddd.QueryEventHandler, options ...*RegisterSubscribeOptions) RegisterSubscribe {
	var subs []*ddd.Subscribe
	if subscribes != nil {
		list := *subscribes
		for i, _ := range list {
			subs = append(subs, &list[i])
		}
	}

	opt := NewRegisterSubscribeOptions(options...)
	return &registerSubscribe{
		subscribes:   subs,
		handler:      handler,
		interceptors: opt.interceptors,
	}
}

func (r *registerSubscribe) GetInterceptor() []ddd.SubscribeInterceptorFunc {
	return r.interceptors
}

func (r *registerSubscribe) GetSubscribes() []*ddd.Subscribe {
	return r.subscribes
}

func (r *registerSubscribe) GetHandler() ddd.QueryEventHandler {
	return r.handler
}

func (r *RegisterEventType) GetEventType() string {
	return r.EventType
}

func (r *RegisterEventType) GetVersion() string {
	return r.Version
}

func (r *RegisterEventType) GetNewFunc() ddd.NewEventFunc {
	return r.NewFunc
}

func RegisterActor(actorServer actor.Server) {
	_actorsFactory = append(_actorsFactory, func() actor.Server { return actorServer })
}

func GetActors() []actor.Factory {
	return _actorsFactory
}

func aggregateSnapshotActorFactory() actor.Server {
	client, err := daprclient.GetDaprDDDClient().DaprClient()
	if err != nil {
		panic(err)
	}
	return ddd.NewAggregateSnapshotActorService(client)
}

func RunWithConfig(setEnv string, configFile string, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventsFunc func() []RegisterEventType, actorsFunc func() []actor.Factory,
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

func RubWithEnvConfig(config *EnvConfig, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventsFunc func() []RegisterEventType, actorsFunc func() []actor.Factory, options ...*RunOptions) (common.Service, error) {
	if len(config.Mongo) > 0 {
		InitMongo(config.App.AppId, config.Mongo)
	}

	if len(config.Neo4j) > 0 {
		InitNeo4j(config.Neo4j)
	}

	if len(config.Minio) > 0 {
		if err := InitMinio(config.Minio); err != nil {
			return nil, err
		}
	}

	opt := NewRunOptions(options...)
	if opt.GetInit() {
		var err error
		if opt.tables != nil {
			err = InitDb(opt.tables, config, opt.GetPrefix())
		}
		return nil, err
	}
	fmt.Printf("---------- %s ----------\r\n", config.App.AppId)
	//创建dapr客户端
	daprClient, err := daprclient.NewDaprDddClient(config.Dapr.GetHost(), config.Dapr.GetHttpPort(), config.Dapr.GetGrpcPort())
	if err != nil {
		panic(err)
	}

	daprclient.SetDaprDddClient(daprClient)

	runCfg := &RunConfig{
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

	return Run(runCfg, config.App.RootUrl, subsFunc, controllersFunc, esMap, eventsFunc, actorsFunc, options...)
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
func Run(runCfg *RunConfig, webRootPath string, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventStorages map[string]ddd.EventStorage,
	eventTypesFunc func() []RegisterEventType, actorsFunc func() []actor.Factory,
	runOptions ...*RunOptions) (common.Service, error) {

	opt := NewRunOptions(runOptions...)
	ddd.Init(runCfg.AppId)
	applog.Init(runCfg.DaprClient, runCfg.AppId, runCfg.LogLevel)
	level := runCfg.LogLevel
	if opt.level != nil {
		level = *opt.level
	}
	serverOptions := &ServiceOptions{
		AppId:          runCfg.AppId,
		HttpHost:       runCfg.HttpHost,
		HttpPort:       runCfg.HttpPort,
		LogLevel:       level,
		EventTypes:     eventTypesFunc(),
		EventStorages:  eventStorages,
		Subscribes:     subsFunc(),
		Controllers:    controllersFunc(),
		ActorFactories: actorsFunc(),
		AuthToken:      "",
		WebRootPath:    webRootPath,
	}
	service := NewService(runCfg.DaprClient, serverOptions)
	if err := service.Start(); err != nil {
		return service, err
	}
	return service, nil
}

func (o *RunOptions) SetFlag(flag *RunFlag) *RunOptions {
	o.SetPrefix(flag.Prefix).SetInit(flag.Init).SetLevel(flag.Level)
	return o
}

func (o *RunOptions) GetInit() bool {
	if o.initDb == nil {
		return false
	}
	return *o.initDb
}

func (o *RunOptions) SetInit(v bool) *RunOptions {
	o.initDb = &v
	return o
}

func (o *RunOptions) GetLevel() *logs.Level {
	return o.level
}

func (o *RunOptions) SetLevel(v *logs.Level) *RunOptions {
	o.level = v
	return o
}

func (o *RunOptions) SetPrefix(v string) *RunOptions {
	o.prefix = &v
	return o
}

func (o *RunOptions) GetPrefix() string {
	if o.prefix == nil {
		return ""
	}
	return *o.prefix
}

func (o *RunOptions) SetTable(v *Tables) *RunOptions {
	o.tables = v
	return o
}

func (o *RunOptions) GetTable() *Tables {
	return o.tables
}
