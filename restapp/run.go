package restapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/setting"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	"github.com/liuxd6825/dapr-go-sdk/service/common"
	"runtime"
	"runtime/debug"
	"strings"
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

type RunOptions struct {
	tables  *Tables
	init    *bool
	sqlFile *string
	prefix  *string
	dbKey   *string
	level   *logs.Level
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

var _currentEnvConfig *EnvConfig

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
		if item.init != nil {
			o.init = item.init
		}
		if item.prefix != nil {
			o.prefix = item.prefix
		}
		if item.dbKey != nil {
			o.dbKey = item.dbKey
		}
		if item.sqlFile != nil {
			o.sqlFile = item.sqlFile
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

func NewRegisterSubscribe(subscribes []*ddd.Subscribe, handler ddd.QueryEventHandler, options ...*RegisterSubscribeOptions) RegisterSubscribe {
	opt := NewRegisterSubscribeOptions(options...)
	return &registerSubscribe{
		subscribes:   subscribes,
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
	service := ddd.NewAggregateSnapshotActorService(client)
	return service
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

	if config == nil {
		return nil, errors.New("config is nil")
	}

	setCpuMemory(config.Name, &config.App)

	if len(config.Mongo) > 0 {
		initMongo(config.App.AppId, config.Mongo)
	}

	if len(config.Neo4j) > 0 {
		initNeo4j(config.Neo4j)
	}

	if len(config.Minio) > 0 {
		if err := initMinio(config.Minio); err != nil {
			return nil, err
		}
	}

	if len(config.Redis) > 0 {
		if err := initRedis(config.Redis); err != nil {
			return nil, err
		}
	}

	if config.App.AuthToken != "" {
		DefaultAuthToken = config.App.AuthToken
	}

	if config.App.AuthTokenKey != "" {
		DefaultAuthTokenKey = config.App.AuthTokenKey
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

	// 启动服务，创建dapr客户端
	daprClient, err := daprclient.NewDaprDddClient(context.Background(), config.Dapr.GetHost(), config.Dapr.GetHttpPort(), config.Dapr.GetGrpcPort())
	if err != nil {
		panic(err)
	}

	daprclient.SetDaprDddClient(daprClient)

	appHost := config.App.HttpHost
	if len(appHost) == 0 {
		appHost = "0.0.0.0"
	}
	runCfg := &RunConfig{
		AppId:                  config.App.AppId,
		HttpHost:               appHost,
		HttpPort:               config.App.HttpPort,
		LogLevel:               config.Log.level,
		DaprMaxCallRecvMsgSize: config.Dapr.MaxCallRecvMsgSize,
		DaprClient:             daprClient,
		EnvConfig:              config,
	}

	eventStoresMap := newEventStores(&config.Dapr, daprClient)

	fmt.Printf("---------- %s ----------\r\n", config.App.AppId)
	return run(runCfg, config.App.RootUrl, subsFunc, controllersFunc, eventStoresMap, eventsFunc, actorsFunc, options...)
}

func newEventStores(cfg *DaprConfig, client daprclient.DaprDddClient) map[string]ddd.EventStore {
	//创建dapr事件存储器
	eventStoresMap := make(map[string]ddd.EventStore)
	esMap := cfg.EventStores
	if len(esMap) == 0 {
		logs.Errorf(context.Background(), "", nil, "config eventStores is empity")
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

// setCpuMemory
//
//	@Description: 设置Cpu和内存大小
//	@param config
func setCpuMemory(envName string, config *AppConfig) {
	if config == nil {
		return
	}
	var fields logs.Fields
	cpu := config.CPU
	maxCpu := runtime.NumCPU()
	if cpu < 0 {
		cpu = maxCpu - cpu
	}
	if maxCpu < cpu {
		cpu = maxCpu
	}
	if cpu <= 0 {
		cpu = 1
	}
	runtime.GOMAXPROCS(cpu)

	memTxt := strings.ToLower(strings.Trim(config.Memory, " "))
	if memTxt == "" {
		logs.Infof(context.Background(), "", fields, "ctype=app; cpu=%v;", cpu)
		return
	}
	var memSize int64 = 0
	size := len(memTxt)
	unit := memTxt[size-1 : size]
	memVal := memTxt[0 : size-1]
	memSize, err := stringutils.ToInt64(memVal)
	if err != nil {
		logs.Panic(context.Background(), "", fields, "ctype=app; memory=%s; 值不正确。示例: 10G, 10M, 10K", envName, memTxt)
	}

	switch unit {
	case "g":
		memSize = memSize * 1024 * 1024 * 1024
	case "m":
		memSize = memSize * 1024 * 1024
	case "k":
		memSize = memSize * 1024
	default:
		logs.Panic(context.Background(), "", fields, "ctype=app; %s.app.memory=%s 不正确。示例: 10G, 10M, 10K", envName, memTxt)
	}
	debug.SetMemoryLimit(memSize)
	logs.Infof(context.Background(), "", fields, "ctype=app; cpu=%v; memory=%s;", cpu, memTxt)
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
		EnvConfig:      runCfg.EnvConfig,
	}

	// 设置全局时区为本地时区
	setting.SetLocalTimeZone()

	_currentEnvConfig = runCfg.EnvConfig

	// 启动HTTP服务器
	service := NewHttpServer(runCfg.DaprClient, serverOptions)
	if err := service.Start(); err != nil {
		return service, err
	}
	return service, nil
}

func (o *RunOptions) SetFlag(flag *RunFlag) *RunOptions {
	o.SetPrefix(flag.Prefix)
	o.SetInit(flag.Init)
	o.SetDbKey(flag.DbKey)
	o.SetSqlFile(flag.SqlFile)
	o.SetLevel(flag.Level)
	return o
}

func (o *RunOptions) GetInit() bool {
	if o.init == nil {
		return false
	}
	return *o.init
}

func (o *RunOptions) SetInit(v bool) *RunOptions {
	o.init = &v
	return o
}

func (o *RunOptions) SetSqlFile(v string) *RunOptions {
	o.sqlFile = &v
	return o
}

func (o *RunOptions) GetSqlFile() string {
	if o.sqlFile == nil {
		return ""
	}
	return *o.sqlFile
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

func (o *RunOptions) SetDbKey(v string) *RunOptions {
	o.dbKey = &v
	return o
}

func (o *RunOptions) GetDbKey() string {
	if o.dbKey == nil {
		return ""
	}
	return *o.dbKey
}

func GetConfigAppValue(name string) (string, error) {
	var err error
	v, ok := _currentEnvConfig.App.Values[name]
	if !ok {
		err = errors.New(fmt.Sprintf("配置变量%s不存在", name))
	}
	return v, err
}

func GetConfigAppValues() map[string]string {
	return _currentEnvConfig.App.Values
}

func SetCurrentEnvConfig(envConfig *EnvConfig) {
	_currentEnvConfig = envConfig
}

func GetCurrentEnvConfig() *EnvConfig {
	return _currentEnvConfig
}
