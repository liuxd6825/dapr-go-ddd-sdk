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

var _app *iris.Application

type StartOptions struct {
	AppId      string
	AppPort    int
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

type RegisterEvent struct {
	EventType string
	Revision  string
	NewFunc   ddd.NewEventFunc
}

func RunWithConfig(configFile string, app *iris.Application, subsFunc func() *[]RegisterSubscribe, controllersFunc func() *[]Controller, eventsFunc func() *[]RegisterEvent) error {
	config, err := NewConfigByFile(configFile)
	if err != nil {
		panic(err)
	}
	envConfig, err := config.GetEnvConfig()
	if err != nil {
		panic(err)
	}
	return RubWithEnvConfig(envConfig, app, subsFunc, controllersFunc, eventsFunc)
}

func RubWithEnvConfig(config *EnvConfig, app *iris.Application, subsFunc func() *[]RegisterSubscribe, controllersFunc func() *[]Controller, eventsFunc func() *[]RegisterEvent) error {
	if !config.Mongo.IsEmpty() {
		initMongo(&config.Mongo)
	}

	//创建dapr客户端
	daprClient, err := daprclient.NewClient(config.Dapr.Host, config.Dapr.HttpPort, config.Dapr.GrpcPort)
	if err != nil {
		panic(err)
	}

	options := &StartOptions{
		AppId:      config.App.AppId,
		AppPort:    config.App.AppPort,
		LogLevel:   config.Log.GetLevel(),
		DaprClient: daprClient,
	}

	eventStorages := map[string]ddd.EventStorage{}

	//创建dapr事件存储器
	esMap := map[string]ddd.EventStorage{}
	for _, pubsubName := range config.Dapr.Pubsubs {
		eventStorage, err := ddd.NewDaprEventStorage(daprClient, ddd.PubsubName(pubsubName))
		if err != nil {
			panic(err)
		}
		esMap[pubsubName] = eventStorage
	}

	return Run(options, app, config.App.RootUrl, subsFunc, controllersFunc, &eventStorages, eventsFunc)
}

//
// Run
// @Description: 启动 iris web 服务
// @param port
// @param app
// @return *iris.Application
// @return error
//
func Run(options *StartOptions, app *iris.Application, rootUrl string, subsFunc func() *[]RegisterSubscribe,
	controllersFunc func() *[]Controller, eventStorages *map[string]ddd.EventStorage, eventsFunc func() *[]RegisterEvent) error {
	_app = app

	ddd.Init(options.AppId)
	applog.Init(options.DaprClient, options.AppId, options.LogLevel)

	// 注册消息订阅
	subs := subsFunc()
	if subs != nil {
		for _, s := range *subs {
			registerSubscribeHandler(s.GetSubscribes(), s.GetHandler())
		}
	}

	// 注册控制器
	controllers := controllersFunc()
	if controllers != nil {
		for _, c := range *controllers {
			registerRestController(rootUrl, c)
		}
	}

	// 注册领域事件类型
	events := eventsFunc()
	if events != nil {
		for _, event := range *events {
			if err := ddd.RegisterEventType(event.EventType, event.Revision, event.NewFunc); err != nil {
				return errors.New(fmt.Sprintf("RegisterEventType() error:%s , EventType=%s, Revision=%s", err.Error(), event.EventType, event.Revision))
			}
		}
	}

	// dapr服务通过访问http://locahost:<port>/dapr/subscribe获取订阅的消息
	_app.Get(subscribePath, func(context *context.Context) {
		_, _ = context.JSON(ddd.GetSubscribes())
	})

	// 注册事件存储器
	if eventStorages != nil {
		for _, es := range *eventStorages {
			ddd.RegisterEventStorage(es.GetPubsubName(), es)
		}
	}
	if err := ddd.StartSubscribeHandlers(); err != nil {
		return err
	}
	if err := app.Run(iris.Addr(fmt.Sprintf(":%d", options.AppPort))); err != nil {
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
func registerSubscribeHandler(subscribes *[]ddd.Subscribe, queryEventHandler ddd.QueryEventHandler) ddd.SubscribeHandler {
	return ddd.NewSubscribeHandler(subscribes, queryEventHandler, func(sh ddd.SubscribeHandler, subscribe ddd.Subscribe) (err error) {
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
}
