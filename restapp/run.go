package restapp

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
)

const subscribePath = "dapr/subscribe"
const eventTypesPath = "dapr/event-types"

var _app *iris.Application
var _webRootPath string

type StartOptions struct {
	AppId      string
	HttpHost   string
	HttpPort   int
	LogLevel   applog.Level
	DaprClient daprclient.DaprClient
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

var _daprClient daprclient.DaprClient

func GetDaprClient() daprclient.DaprClient {
	return _daprClient
}

func RunWithConfig(envType string, configFile string, app *iris.Application, subsFunc func() *[]RegisterSubscribe, controllersFunc func() *[]Controller, eventsFunc func() *[]RegisterEventType) error {
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
	return RubWithEnvConfig(envConfig, app, subsFunc, controllersFunc, eventsFunc)
}

func RubWithEnvConfig(config *EnvConfig, app *iris.Application, subsFunc func() *[]RegisterSubscribe, controllersFunc func() *[]Controller, eventsFunc func() *[]RegisterEventType) error {
	if !config.Mongo.IsEmpty() {
		initMongo(&config.Mongo)
	}

	//创建dapr客户端
	daprClient, err := daprclient.NewClient(config.Dapr.GetHost(), config.Dapr.GetHttpPort(), config.Dapr.GetGrpcPort())
	if err != nil {
		panic(err)
	}

	_daprClient = daprClient

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

	return Run(options, app, config.App.RootUrl, subsFunc, controllersFunc, esMap, eventsFunc)
}

//
// Run
// @Description:
// @param options
// @param app
// @param webRootPath web server URL root path
// @param subsFunc
// @param controllersFunc
// @param eventStorages
// @param eventTypesFunc
// @return error
//
func Run(options *StartOptions, app *iris.Application, webRootPath string, subsFunc func() *[]RegisterSubscribe,
	controllersFunc func() *[]Controller, eventStorages map[string]ddd.EventStorage, eventTypesFunc func() *[]RegisterEventType) error {
	_app = app
	_webRootPath = webRootPath
	ddd.Init(options.AppId)
	applog.Init(options.DaprClient, options.AppId, options.LogLevel)

	// 注册消息订阅
	subs := subsFunc()
	if subs != nil {
		for _, s := range *subs {
			if s != nil {
				if _, err := registerSubscribeHandler(s.GetSubscribes(), s.GetHandler()); err != nil {
					return err
				}
			}
		}
	}

	// 注册控制器
	controllers := controllersFunc()
	if controllers != nil {
		for _, c := range *controllers {
			if c != nil {
				registerRestController(webRootPath, c)
			}
		}
	}

	// 注册领域事件类型
	eventTypes := eventTypesFunc()
	if eventTypes != nil {
		for _, t := range *eventTypes {
			if err := ddd.RegisterEventType(t.EventType, t.Revision, t.NewFunc); err != nil {
				return errors.New(fmt.Sprintf("RegisterEventType() error:%s , EventType=%s, Revision=%s", err.Error(), t.EventType, t.Revision))
			}
		}
	}

	// dapr服务通过访问http://locahost:<port>/dapr/subscribe获取订阅的消息
	_app.Get(subscribePath, func(context *context.Context) {
		data := ddd.GetSubscribes()
		_, _ = context.JSON(data)
	})

	// dapr服务通过访问http://locahost:<port>/dapr/subscribe获取订阅的消息
	_app.Get(eventTypesPath, func(context *context.Context) {
		//_, _ = context.JSON(ddd.RegisterEventType())
	})

	// 注册事件存储器
	if eventStorages != nil {
		for key, es := range eventStorages {
			ddd.RegisterEventStorage(key, es)
		}
	}
	if err := ddd.StartSubscribeHandlers(); err != nil {
		return err
	}
	if err := app.Run(iris.Addr(fmt.Sprintf("%s:%d", options.HttpHost, options.HttpPort))); err != nil {
		return err
	}
	return nil
}

func NewRegisterController(relativePath string, ctls ...interface{}) RegisterController {
	return RegisterController{
		RelativePath: relativePath,
		Controllers:  ctls,
	}
}

//
// RegisterRestController
// @Description: 注册UserInterface层Controller
// @param relativePath
// @param configurators
//
func registerRestController(relativePath string, controllers ...Controller) {
	if controllers == nil && len(controllers) == 0 {
		return
	}
	configurators := func(app *mvc.Application) {
		for _, c := range controllers {
			app.Handle(c)
		}
	}
	mvc.Configure(_app.Party(relativePath), configurators)
}

//
// registerQueryHandler
// @Description: 注册领域事件控制器
// @param handlers
// @return error
//
func registerQueryHandler(handlers ...ddd.SubscribeHandler) error {
	// 注册User消息处理器
	for _, h := range handlers {
		err := ddd.RegisterQueryHandler(h)
		if err != nil {
			return err
		}
	}
	return nil
}

//
// registerSubscribeHandler
// @Description: 新建领域事件控制器
// @param subscribes
// @param queryEventHandler
// @return ddd.SubscribeHandler
//
func registerSubscribeHandler(subscribes *[]ddd.Subscribe, queryEventHandler ddd.QueryEventHandler) (ddd.SubscribeHandler, error) {
	handler := ddd.NewSubscribeHandler(subscribes, queryEventHandler, func(sh ddd.SubscribeHandler, subscribe ddd.Subscribe) (err error) {
		defer func() {
			if e := ddd_errors.GetRecoverError(recover()); e != nil {
				err = e
			}
		}()
		_app.Handle("POST", subscribe.Route, func(c *context.Context) {
			if err = sh.CallQueryEventHandler(c, c); err != nil {
				c.SetErr(err)
			}
		})
		return err
	})
	if err := ddd.RegisterQueryHandler(handler); err != nil {
		return nil, err
	}
	return handler, nil
}
