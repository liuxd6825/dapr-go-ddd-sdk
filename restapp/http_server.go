package restapp

import (
	context2 "context"
	"fmt"
	"github.com/dapr/go-sdk/actor"
	"github.com/dapr/go-sdk/actor/runtime"
	"github.com/dapr/go-sdk/service/common"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/html2/template"
	swagger3 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/swagger/v3"
	"net/http"
	"time"
)

type ServiceOptions struct {
	AppId          string
	HttpHost       string
	HttpPort       int
	LogLevel       applog.Level
	ActorFactories []actor.FactoryContext
	Subscribes     []RegisterSubscribe
	Controllers    []Controller
	EventTypes     []RegisterEventType
	AuthToken      string
	WebRootPath    string
	SwaggerDoc     string
	EnvConfig      *EnvConfig
	OnInitEvents   []OnInitEvent
	OnStartEvents  []OnStartEvent
}

type HttpServer struct {
	app              *iris.Application
	appId            string
	httpHost         string
	httpPort         int
	logLevel         applog.Level
	daprDddClient    dapr.DaprClient
	eventStores      map[string]ddd.EventStore
	actorFactories   []actor.FactoryContext
	subscribes       []RegisterSubscribe
	controllers      []Controller
	eventTypes       []RegisterEventType
	jobEventHandlers map[string]common.JobEventHandler
	authToken        string
	webRootPath      string
	envConfig        *EnvConfig
	onInitEvents     []OnInitEvent
	onStartEvents    []OnStartEvent
	swagger          *swagger3.Swagger
	fs               *fsm.Manager

	httpServer *http.Server
	daprServer common.Service
}

type OnStartEvent func(server *HttpServer) error

var _startEvents []OnStartEvent

func RegisterOnStartInit(startEvent OnStartEvent) {
	if startEvent != nil {
		_startEvents = append(_startEvents, startEvent)
	}
}

func NewHttpServer(daprDddClient dapr.DaprClient, opts *ServiceOptions) common.Service {
	actorRuntime := runtime.GetActorRuntimeInstanceContext()
	envConfig := opts.EnvConfig

	if opts.EnvConfig != nil {
		actorConfig := actorRuntime.Config()
		actorConfig.DrainOngingCallTimeout = envConfig.Dapr.Actor.DrainOngingCallTimeout
		actorConfig.ActorScanInterval = envConfig.Dapr.Actor.ActorScanInterval
		actorConfig.ActorIdleTimeout = envConfig.Dapr.Actor.ActorIdleTimeout
		actorConfig.DrainBalancedActors = envConfig.Dapr.Actor.DrainBalancedActors
	}

	return &HttpServer{
		app:              iris.New(),
		httpPort:         opts.HttpPort,
		httpHost:         opts.HttpHost,
		appId:            opts.AppId,
		logLevel:         opts.LogLevel,
		daprDddClient:    daprDddClient,
		actorFactories:   opts.ActorFactories,
		subscribes:       opts.Subscribes,
		controllers:      opts.Controllers,
		eventTypes:       opts.EventTypes,
		authToken:        opts.AuthToken,
		webRootPath:      opts.WebRootPath,
		envConfig:        opts.EnvConfig,
		onInitEvents:     opts.OnInitEvents,
		onStartEvents:    opts.OnStartEvents,
		swagger:          swagger3.NewSwagger(),
		jobEventHandlers: make(map[string]common.JobEventHandler),
	}

}

// EnvConfig
//
//	@Description:
//	@receiver s
//	@return *EnvConfig
func (s *HttpServer) EnvConfig() *EnvConfig {
	return s.envConfig
}

// App
//
//	@Description:
//	@receiver s
//	@return *iris.Application
func (s *HttpServer) App() *iris.Application {
	return s.app
}

// DaprClient
//
//	@Description:
//	@receiver s
//	@return dapr.DaprClient
func (s *HttpServer) DaprClient() dapr.DaprClient {
	return s.daprDddClient
}

// Start
//
//	@Description:
//	@Description:
//	@receiver s
//	@return error
func (s *HttpServer) Start() error {
	ctx := logs.NewContext(context2.Background())
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			logs.Info(ctx, "", logs.Fields{"func": "restapp.HttpServer.Start()", "error": err.Error()})
		}
	}()
	app := s.app

	// 使用自定义的 JSON 编码器替换 Iris 默认的 JSON 编码器
	app.Configure(iris.WithOptimizations)

	/*
		if err := s.addRenderHandler(app); err != nil {
			return err
		}
	*/

	// 注册控制器
	if s.controllers != nil {
		for _, c := range s.controllers {
			if c != nil {
				s.registerController(s.webRootPath, c)
			}
		}
	}

	s.registerBaseHandler()
	for _, init := range s.onInitEvents {
		if err := init(s); err != nil {
			return err
		}
	}

	if err := s.addSwaggerHandler(app); err != nil {
		panic(err.Error())
	}

	var actors []actor.FactoryContext
	actors = append(actors, s.actorFactories...)
	actors = append(actors, GetActors()...)
	for _, f := range actors {
		s.RegisterActorImplFactoryContext(f)
	}

	app.ConfigureHost(func(su *host.Supervisor) {
		// httpServer:=su.Server
		// println("httpServer", httpServer)
	})

	addr := fmt.Sprintf("%s:%d", s.httpHost, s.httpPort)
	if err := app.Run(iris.Addr(addr), func(app *iris.Application) {
		if err := s.doOnStartEvents(ctx, app); err != nil {
			panic(err)
		}
	}); err != nil {
		return err
	}

	return nil
}

func (s *HttpServer) doOnStartEvents(ctx context2.Context, app *iris.Application) error {
	if err := s.startSubscribeHandlers(); err != nil {
		panic(err.Error())
	}

	fmt.Printf("---------- %s running ----------\r\n", s.envConfig.App.AppId)
	if logs.GetLevel() <= logs.DebugLevel {
		for _, v := range app.GetRoutes() {
			logs.Debug(ctx, "", logs.Fields{"route": v.Method + " " + v.Path})
		}
	}

	var startEvents []OnStartEvent
	startEvents = append(startEvents, _startEvents...)
	startEvents = append(startEvents, s.onStartEvents...)

	for _, event := range startEvents {
		if err := event(s); err != nil {
			return err
		}
	}
	return nil
}

func (s *HttpServer) startSubscribeHandlers() error {
	// 注册消息订阅
	if s.subscribes != nil {
		for _, subscribe := range s.subscribes {
			if subscribe != nil {
				if _, err := s.registerSubscribeHandler(subscribe.GetSubscribes(), subscribe.GetHandler(), subscribe.GetInterceptor()); err != nil {
					panic(err.Error())
				}
			}
		}
	}
	if err := ddd.StartSubscribeHandlers(); err != nil {
		return err
	}
	return nil
}

func (s *HttpServer) addSwaggerHandler(app *iris.Application) error {
	path := "swagger.json"
	app.Get(path, func(ictx iris.Context) {
		resource := ictx.URLParam("resource")
		swagger := s.swagger.Filter(resource)
		if err := ictx.JSON(swagger); err != nil {
			ictx.StatusCode(http.StatusInternalServerError)
			ictx.SetErr(err)
		}
	})
	return nil
}

func (s *HttpServer) addRenderHandler(app *iris.Application) error {
	if !s.envConfig.App.Template.Enable {
		return nil
	}

	fsKey := s.envConfig.App.Template.FsId
	if fsKey == "" {
		return errors.New("template fsKey is empty")
	}

	fs, ok := s.envConfig.fsManager.GetFs(fsKey)
	if !ok {
		return errors.New("template fsManager fs key not found %s", fsKey)
	}

	render, err := template.NewHandler(fs)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/{filePath:path}", s.envConfig.App.Template.ApiUrl)
	app.Get(apiUrl, func(ictx iris.Context) {
		filePath := ictx.Params().Get("filePath")
		render(ictx, filePath)
	})
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

func (s *HttpServer) healthHandler(context *context.Context) {
	context.StatusCode(http.StatusOK)
}

func (s *HttpServer) eventTypesHandler(ctx *context.Context) {

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
		err := ddd.RegisterQueryHandler(h, ddd.GetEventStoreDefaultPubsubName())
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *HttpServer) GetSwagger() *swagger3.Swagger {
	return s.swagger
}
func (s *HttpServer) AddSwagger(swg *swagger3.Swagger) {
	if swg == nil {
		return
	}
	for key, path := range swg.Paths {
		s.swagger.Paths[key] = path
	}
}
func (s *HttpServer) AddSubscribe(subs ...RegisterSubscribe) {
	s.subscribes = append(s.subscribes, subs...)
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
