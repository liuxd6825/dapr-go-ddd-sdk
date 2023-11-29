package restapp

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
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
}

type RunOptions struct {
	tables   *Tables
	initDb   *bool
	dbScript *bool
	prefix   *string
	dbKey    *string
	file     *string
	dbkey    *string
	level    *logs.Level
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

var _actorsFactory []actor.FactoryContext = []actor.FactoryContext{
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

func RegisterActor(actorServer actor.ServerContext) {
	_actorsFactory = append(_actorsFactory, func() actor.ServerContext { return actorServer })
}

func GetActors() []actor.FactoryContext {
	return _actorsFactory
}

func aggregateSnapshotActorFactory() actor.ServerContext {
	client, err := daprclient.GetDaprDDDClient().DaprClient()
	if err != nil {
		panic(err)
	}
	return ddd.NewAggregateSnapshotActorService(client)
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

func RubWithEnvConfig(config *EnvConfig, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventsFunc func() []RegisterEventType, actorsFunc func() []actor.FactoryContext, options ...*RunOptions) (common.Service, error) {
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
	if opt.GetInitDb() {
		var err error
		if opt.tables != nil {
			err = InitDb(opt.GetDbKey(), opt.tables, config, opt.GetPrefix())
		}
		return nil, err
	}

	if opt.GetDbScript() {
		var err error
		if opt.tables != nil {
			err = InitDbScript(opt.GetDbKey(), opt.tables, config, opt.GetPrefix())
		}
		return nil, err
	}

	//创建dapr客户端
	daprClient, err := daprclient.NewDaprDddClient(context.Background(), config.Dapr.GetHost(), config.Dapr.GetHttpPort(), config.Dapr.GetGrpcPort())
	if err != nil {
		panic(err)
	}

	daprclient.SetDaprDddClient(daprClient)

	runCfg := &RunConfig{
		AppId:                  config.App.AppId,
		HttpHost:               config.App.HttpHost,
		HttpPort:               config.App.HttpPort,
		LogLevel:               config.Log.GetLevel(),
		DaprMaxCallRecvMsgSize: config.Dapr.MaxCallRecvMsgSize,
		DaprClient:             daprClient,
	}

	eventStoresMap := newEventStores(&config.Dapr, daprClient)

	fmt.Printf("---------- %s ----------\r\n", config.App.AppId)
	return Run(runCfg, config.App.RootUrl, subsFunc, controllersFunc, eventStoresMap, eventsFunc, actorsFunc, options...)
}

func newEventStores(cfg *DaprConfig, client daprclient.DaprDddClient) map[string]ddd.EventStore {
	//创建dapr事件存储器
	eventStoresMap := make(map[string]ddd.EventStore)
	esMap := cfg.EventStores
	if len(esMap) == 0 {
		panic("config eventStores is empity")
	} else {
		var defEs ddd.EventStore
		for _, item := range esMap {
			eventStorage, err := ddd.NewGrpcEventStore(item.CompName, item.PubsubName, client)
			if err != nil {
				panic(err)
			}
			eventStoresMap[item.CompName] = eventStorage
			if defEs == nil {
				defEs = eventStorage
			}
		}
		eventStoresMap[""] = defEs
	}
	return eventStoresMap
}

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
func Run(runCfg *RunConfig, webRootPath string, subsFunc func() []RegisterSubscribe,
	controllersFunc func() []Controller, eventStores map[string]ddd.EventStore,
	eventTypesFunc func() []RegisterEventType, actorsFunc func() []actor.FactoryContext,
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
		EventStores:    eventStores,
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
	o.SetPrefix(flag.Prefix).SetInitDb(flag.InitDb).SetLevel(flag.Level)
	return o
}

func (o *RunOptions) GetInitDb() bool {
	if o.initDb == nil {
		return false
	}
	return *o.initDb
}

func (o *RunOptions) GetDbScript() bool {
	if o.dbScript == nil {
		return false
	}
	return *o.dbScript
}

func (o *RunOptions) SetInitDb(v bool) *RunOptions {
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

func (o *RunOptions) SetFile(v *string) *RunOptions {
	o.file = v
	return o
}

func (o *RunOptions) GetFile() string {
	return *o.file
}

func (o *RunOptions) SetDbKey(v *string) *RunOptions {
	o.dbkey = v
	return o
}

func (o *RunOptions) GetDbKey() string {
	if o.dbkey == nil {
		return ""
	}
	return *o.dbkey
}
