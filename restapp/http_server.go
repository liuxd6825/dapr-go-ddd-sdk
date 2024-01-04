package restapp

import (
	context2 "context"
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
	"github.com/liuxd6825/dapr-go-sdk/actor/runtime"
	"github.com/liuxd6825/dapr-go-sdk/service/common"
	"net/http"
	"strings"
	"time"
)

type ServiceOptions struct {
	AppId          string
	HttpHost       string
	HttpPort       int
	LogLevel       applog.Level
	EventStores    map[string]ddd.EventStore
	ActorFactories []actor.FactoryContext
	Subscribes     []RegisterSubscribe
	Controllers    []Controller
	EventTypes     []RegisterEventType
	AuthToken      string
	WebRootPath    string
	SwaggerDoc     string
	EnvConfig      *EnvConfig
}

type HttpServer struct {
	app                         *iris.Application
	appId                       string
	httpHost                    string
	httpPort                    int
	logLevel                    applog.Level
	daprDddClient               daprclient.DaprDddClient
	eventStores                 map[string]ddd.EventStore
	actorFactories              []actor.FactoryContext
	subscribes                  []RegisterSubscribe
	controllers                 []Controller
	eventTypes                  []RegisterEventType
	authToken                   string
	webRootPath                 string
	eventStoreDefaultPubsubName string // 默认事件存储器的名称

	sdkServer *http.Server
}

type OnAppInit func(ctx context2.Context) error

var _appInits []OnAppInit

func RegisterOnAppInit(init OnAppInit) {
	if init != nil {
		_appInits = append(_appInits, init)
	}
}

func NewHttpServer(daprDddClient daprclient.DaprDddClient, opts *ServiceOptions) common.Service {
	eventStoreDefaultPubsubName := ""
	es, ok := opts.EventStores[""]
	if ok {
		eventStoreDefaultPubsubName = es.GetPubsubName()
	}

	actorRuntime := runtime.GetActorRuntimeInstanceContext()
	envConfig := opts.EnvConfig
	if opts.EnvConfig != nil {
		actorConfig := actorRuntime.Config()
		actorConfig.DrainOngingCallTimeout = envConfig.Dapr.Actor.DrainOngingCallTimeout
		actorConfig.ActorScanInterval = envConfig.Dapr.Actor.ActorScanInterval
		actorConfig.ActorIdleTimeout = envConfig.Dapr.Actor.ActorIdleTimeout
		actorConfig.DrainBalancedActors = envConfig.Dapr.Actor.DrainBalancedActors
	}

	app := iris.New()
	return &HttpServer{
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
		app:                         app,
	}

}

func (s *HttpServer) Start() error {
	ctx := logs.NewContext(context2.Background())
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			logs.Info(ctx, "", logs.Fields{"func": "restapp.HttpServer.Start()", "error": err.Error()})
		}
	}()
	app := s.app

	s.registerBaseHandler()

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
			s.RegisterActorImplFactoryContext(f)
		}
	}

	addr := fmt.Sprintf("%s:%d", s.httpHost, s.httpPort)
	if err := app.Run(iris.Addr(addr), func(application *iris.Application) {
		for _, onInit := range _appInits {
			if err := onInit(ctx); err != nil {
				panic(err.Error())
			}
		}
	}); err != nil {
		return err
	}

	return nil
}

// AddHealthCheckHandler appends provided app health check handler.
func (s *HttpServer) AddHealthCheckHandler(route string, fn common.HealthCheckHandler) error {
	if fn == nil {
		return fmt.Errorf("health check handler required")
	}

	if !strings.HasPrefix(route, "/") {
		route = fmt.Sprintf("/%s", route)
	}

	s.app.HandleMany("ALL", route, optionsHandler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := fn(r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})))
	return nil
}

func (s *HttpServer) Stop() error {
	ctxShutDown, cancel := context2.WithTimeout(context2.Background(), 5*time.Second)
	defer cancel()

	return s.app.Shutdown(ctxShutDown)
}

func (s *HttpServer) GracefulStop() error {
	return s.Stop()
}

func (s *HttpServer) subscribesHandler(ictx *context.Context) {
	ctx, _ := NewContext(ictx)
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			fields := logs.Fields{
				"func":  "restapp.HttpServer.subscribesHandler()",
				"error": err.Error(),
			}
			logs.Error(ctx, "", fields)
		}
	}()

	subscribes := ddd.GetSubscribes()
	ictx.Header("Context-Type", "application/json")
	_ = ictx.JSON(subscribes)

	if logs.GetLevel(ctx) >= logs.InfoLevel {
		for _, s := range subscribes {
			fields := logs.Fields{
				"dapr":   "subscribes",
				"pubsub": s.PubsubName,
				"topic":  s.Topic,
				"route":  s.Route,
			}
			logs.Info(ctx, "", fields)
		}
	}

}

func (s *HttpServer) healthHandler(context *context.Context) {
	context.StatusCode(http.StatusOK)
}

func (s *HttpServer) eventTypesHandler(ctx *context.Context) {

}

// registerSubscribeHandler
// @Description: 新建领域事件控制器
// @param subscribes
// @param queryEventHandler
// @return ddd.SubscribeHandler
func (s *HttpServer) registerSubscribeHandler(subscribes []*ddd.Subscribe, queryEventHandler ddd.QueryEventHandler, interceptors []ddd.SubscribeInterceptorFunc) (ddd.SubscribeHandler, error) {
	subscribesHandler := func(sh ddd.SubscribeHandler, subscribe *ddd.Subscribe) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		s.app.Handle("POST", subscribe.Route, func(ictx *context.Context) {
			c, _ := NewContext(ictx)
			if err = sh.SubscribeHandler(c, ictx); err != nil {
				ictx.SetErr(err)
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
func (s *HttpServer) registerController(relativePath string, controllers ...Controller) {
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
func (s *HttpServer) registerSwagger() {
	url := fmt.Sprintf("http://%s:%d/swagger/doc.json", "localhost", s.httpPort)
	cfg := &swagger.Config{
		URL: url,
	}
	// use swagger middleware to
	s.app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(cfg, swaggerFiles.Handler))
}

// registerQueryHandler
// @Description: 注册领域事件控制器
// @param handlers
// @return error
func (s *HttpServer) registerQueryHandler(handlers ...ddd.SubscribeHandler) error {
	// 注册User消息处理器
	for _, h := range handlers {
		err := ddd.RegisterQueryHandler(h, s.eventStoreDefaultPubsubName)
		if err != nil {
			return err
		}
	}
	return nil
}

func setOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	w.Header().Set("Allow", "POST,OPTIONS")
}

func optionsHandler(h http.Handler) context.Handler {
	return func(c *context.Context) {
		if c.Method() == http.MethodOptions {
			setOptions(c.ResponseWriter(), c.Request())
		} else {
			h.ServeHTTP(c.ResponseWriter(), c.Request())
		}
	}
}
