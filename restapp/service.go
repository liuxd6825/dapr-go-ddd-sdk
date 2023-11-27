package restapp

import (
	"fmt"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	"github.com/liuxd6825/dapr-go-sdk/actor/config"
	actorErr "github.com/liuxd6825/dapr-go-sdk/actor/error"
	"github.com/liuxd6825/dapr-go-sdk/actor/runtime"
	"github.com/liuxd6825/dapr-go-sdk/service/common"
	"net/http"
)

type ServiceOptions struct {
	AppId          string
	HttpHost       string
	HttpPort       int
	LogLevel       applog.Level
	EventStores    map[string]ddd.EventStore
	ActorFactories []actor.Factory
	Subscribes     []RegisterSubscribe
	Controllers    []Controller
	EventTypes     []RegisterEventType
	AuthToken      string
	WebRootPath    string
	SwaggerDoc     string
}

type service struct {
	app                         *iris.Application
	appId                       string
	httpHost                    string
	httpPort                    int
	logLevel                    applog.Level
	daprDddClient               daprclient.DaprDddClient
	eventStores                 map[string]ddd.EventStore
	actorFactories              []actor.Factory
	subscribes                  []RegisterSubscribe
	controllers                 []Controller
	eventTypes                  []RegisterEventType
	authToken                   string
	webRootPath                 string
	eventStoreDefaultPubsubName string // 默认事件存储器的名称
}

func (s *service) AddHealthCheckHandler(name string, fn common.HealthCheckHandler) error {
	return nil
}

func (s *service) RegisterActorImplFactoryContext(f actor.FactoryContext, opts ...config.Option) {
	runtime.GetActorRuntimeInstance().RegisterActorFactoryContext(f, opts...)
}

func (s *service) RegisterActorImplFactory(f actor.Factory, opts ...config.Option) {
	runtime.GetActorRuntimeInstance().RegisterActorFactory(f, opts...)
}

func NewService(daprDddClient daprclient.DaprDddClient, opts *ServiceOptions) common.Service {
	eventStoreDefaultPubsubName := ""
	es, ok := opts.EventStores[""]
	if ok {
		eventStoreDefaultPubsubName = es.GetPubsubName()
	}

	return &service{
		httpPort:                    opts.HttpPort,
		httpHost:                    opts.HttpHost,
		appId:                       opts.AppId,
		logLevel:                    opts.LogLevel,
		daprDddClient:               daprDddClient,
		eventStores:                 opts.EventStores,
		actorFactories:              opts.ActorFactories,
		subscribes:                  opts.Subscribes,
		controllers:                 opts.Controllers,
		eventTypes:                  opts.EventTypes,
		authToken:                   opts.AuthToken,
		webRootPath:                 opts.WebRootPath,
		eventStoreDefaultPubsubName: eventStoreDefaultPubsubName,
		app:                         iris.New(),
	}
}

func (s *service) AddServiceInvocationHandler(name string, fn common.ServiceInvocationHandler) error {
	return nil
}

func (s *service) AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error {
	return nil
}

func (s *service) AddBindingInvocationHandler(name string, fn common.BindingInvocationHandler) error {
	return nil
}

func (s *service) Stop() error {
	return nil
}

func (s *service) GracefulStop() error {
	return nil
}

func (s *service) setOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	w.Header().Set("Allow", "POST,OPTIONS")
}

func (s *service) Start() error {
	app := s.app

	//app.Use(GlobalJsonSerialization)
	// register subscribe handler
	app.Get("dapr/subscribe", s.subscribesHandler)

	// register domain event types
	app.Get("dapr/event-types", s.eventTypesHandler)

	//	register health check handler
	app.Get("/healthz", s.healthHandler)

	// register actor config handler
	app.Get("/dapr/config", s.actorConfigHandler)

	// register actor method invoke handler
	app.Put("/actors/{actorType}/{actorId}/method/{methodName}", s.actorMethodInvokeHandler)

	// register actor reminder invoke handler
	app.Put("/actors/{actorType}/{actorId}/method/remind/{reminderName}", s.actorReminderInvokeHandler)

	// register actor reminder invoke handler
	app.Put("/actors/{actorType}/{actorId}/method/timer/{timerName}", s.actorTimerInvokeHandler)

	// register deactivate actor handler
	app.Delete("/actors/{actorType}/{actorId}", s.actorDeactivateHandler)

	// register swagger doc
	s.registerSwagger()

	// 注册消息订阅
	if s.subscribes != nil {
		for _, subscribe := range s.subscribes {
			if subscribe != nil {
				if _, err := s.registerSubscribeHandler(subscribe.GetSubscribes(), subscribe.GetHandler(), subscribe.GetInterceptor()); err != nil {
					return err
				}
			}
		}
	}

	// 注册控制器
	if s.controllers != nil {
		for _, c := range s.controllers {
			if c != nil {
				s.registerController(s.webRootPath, c)
			}
		}
	}

	// 注册领域事件类型
	if s.eventTypes != nil {
		for _, t := range s.eventTypes {
			if err := ddd.RegisterEventType(t.EventType, t.Version, t.NewFunc); err != nil {
				return errors.New(fmt.Sprintf("RegisterEventType() error:\"%s\" , EventType=\"%s\", Version=\"%s\"", err.Error(), t.EventType, t.Version))
			}
		}
	}

	// 注册事件存储器
	if s.eventStores != nil {
		for key, es := range s.eventStores {
			ddd.RegisterEventStore(key, es)
		}
	}
	if err := ddd.StartSubscribeHandlers(); err != nil {
		return err
	}

	if s.actorFactories != nil {
		for _, f := range s.actorFactories {
			s.RegisterActorImplFactory(f)
		}
	}
	addr := fmt.Sprintf("%s:%d", s.httpHost, s.httpPort)
	if err := app.Run(iris.Addr(addr)); err != nil {
		return err
	}
	return nil
}

// register actor method invoke handler
func (s *service) actorMethodInvokeHandler(ctx *context.Context) {
	actorType := ctx.Params().Get("actorType")
	actorId := ctx.Params().Get("actorId")
	methodName := ctx.Params().Get("methodName")
	reqData, _ := ctx.GetBody()
	rspData, err := runtime.GetActorRuntimeInstance().InvokeActorMethod(actorType, actorId, methodName, reqData)
	if err == actorErr.ErrActorTypeNotFound {
		ctx.StatusCode(http.StatusNotFound)
		return
	}
	if err != actorErr.Success {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}

	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.Write(rspData)
}

// register actor reminder invoke handler
func (s *service) actorReminderInvokeHandler(ctx *context.Context) {
	actorType := ctx.Params().Get("actorType")
	actorID := ctx.Params().Get("actorId")
	reminderName := ctx.Params().Get("reminderName")
	reqData, _ := ctx.GetBody()
	err := runtime.GetActorRuntimeInstance().InvokeReminder(actorType, actorID, reminderName, reqData)
	if err == actorErr.ErrActorTypeNotFound {
		ctx.StatusCode(http.StatusNotFound)
		return
	}
	if err != actorErr.Success {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
}

// register actor timer invoke handler
func (s *service) actorTimerInvokeHandler(ctx *context.Context) {
	actorType := ctx.Params().Get("actorType")
	actorID := ctx.Params().Get("actorId")
	timerName := ctx.Params().Get("timerName")
	reqData, _ := ctx.GetBody()
	err := runtime.GetActorRuntimeInstance().InvokeTimer(actorType, actorID, timerName, reqData)
	if err == actorErr.ErrActorTypeNotFound {
		ctx.StatusCode(http.StatusNotFound)
		return
	}
	if err != actorErr.Success {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
}

// register deactivate actor handler
func (s *service) actorDeactivateHandler(ctx *context.Context) {
	actorType := ctx.Params().Get("actorType")
	actorID := ctx.Params().Get("actorId")
	err := runtime.GetActorRuntimeInstance().Deactivate(actorType, actorID)
	if err == actorErr.ErrActorTypeNotFound || err == actorErr.ErrActorIDNotFound {
		ctx.StatusCode(http.StatusNotFound)
		return
	}
	if err != actorErr.Success {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
}

func (s *service) subscribesHandler(ictx *context.Context) {
	subscribes := ddd.GetSubscribes()
	_ = ictx.JSON(subscribes)

	ctx := NewContext(ictx)
	for _, s := range subscribes {
		logs.Infof(ctx, "subscribe  pubsubName:%s,  topic:%s,  route:%s,  metadata:%s", s.PubsubName, s.Topic, s.Route, s.Metadata)
	}
}

func (s *service) eventTypesHandler(ctx *context.Context) {

}

func (s *service) healthHandler(context *context.Context) {
	context.StatusCode(http.StatusOK)
}

// register actor config handler
func (s *service) actorConfigHandler(ctx *context.Context) {
	data, err := runtime.GetActorRuntimeInstance().GetJSONSerializedConfig()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.StatusCode(http.StatusOK)
	if _, err = ctx.Write(data); err != nil {
		return
	}
}

// registerQueryHandler
// @Description: 注册领域事件控制器
// @param handlers
// @return error
func (s *service) registerQueryHandler(handlers ...ddd.SubscribeHandler) error {
	// 注册User消息处理器
	for _, h := range handlers {
		err := ddd.RegisterQueryHandler(h, s.eventStoreDefaultPubsubName)
		if err != nil {
			return err
		}
	}
	return nil
}

// registerSubscribeHandler
// @Description: 新建领域事件控制器
// @param subscribes
// @param queryEventHandler
// @return ddd.SubscribeHandler
func (s *service) registerSubscribeHandler(subscribes []*ddd.Subscribe, queryEventHandler ddd.QueryEventHandler, interceptors []ddd.SubscribeInterceptorFunc) (ddd.SubscribeHandler, error) {
	subscribesHandler := func(sh ddd.SubscribeHandler, subscribe *ddd.Subscribe) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		s.app.Handle("POST", subscribe.Route, func(c *context.Context) {
			ctx := logs.NewContext(c, _logger)
			if err = sh.Handler(ctx, c); err != nil {
				c.SetErr(err)
			}
		})
		return err
	}

	handler := ddd.NewSubscribeHandler(subscribes, queryEventHandler, subscribesHandler, interceptors)
	if err := ddd.RegisterQueryHandler(handler, s.eventStoreDefaultPubsubName); err != nil {
		return nil, err
	}
	return handler, nil
}

// RegisterRestController
// @Description: 注册UserInterface层Controller
// @param relativePath
// @param configurators
func (s *service) registerController(relativePath string, controllers ...Controller) {
	if controllers == nil && len(controllers) == 0 {
		return
	}
	configurators := func(app *mvc.Application) {
		for _, c := range controllers {
			app.Handle(c)
		}
	}
	for _, c := range controllers {
		if reg, ok := c.(RegisterHandler); ok {
			reg.RegisterHandler(s.app)
		}
	}
	mvc.Configure(s.app.Party(relativePath), configurators)
}

// registerSwagger
// @Description:
// @receiver s
func (s *service) registerSwagger() {
	url := fmt.Sprintf("http://%s:%d/swagger/doc.json", "localhost", s.httpPort)
	cfg := &swagger.Config{
		URL: url,
	}
	// use swagger middleware to
	s.app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(cfg, swaggerFiles.Handler))
}
