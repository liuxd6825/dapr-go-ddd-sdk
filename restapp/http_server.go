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
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template"
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
	ActorFactories []actor.FactoryContext
	Subscribes     []RegisterSubscribe
	Controllers    []Controller
	EventTypes     []RegisterEventType
	AuthToken      string
	WebRootPath    string
	SwaggerDoc     string
	EnvConfig      *EnvConfig
	Inits          []RunInitFunc
}

type HttpServer struct {
	app            *iris.Application
	appId          string
	httpHost       string
	httpPort       int
	logLevel       applog.Level
	daprDddClient  dapr.DaprClient
	eventStores    map[string]ddd.EventStore
	actorFactories []actor.FactoryContext
	subscribes     []RegisterSubscribe
	controllers    []Controller
	eventTypes     []RegisterEventType
	authToken      string
	webRootPath    string
	sdkServer      *http.Server
	envConfig      *EnvConfig
	inits          []RunInitFunc
}

type OnAppInit func(ctx context2.Context) error

var _appInits []OnAppInit

func RegisterOnAppInit(init OnAppInit) {
	if init != nil {
		_appInits = append(_appInits, init)
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

	app := iris.New()
	//tmpl := iris.HTML("./views", ".html")

	// Enable re-build on local template files changes.
	//tmpl.Reload(true)

	// Register the view engine to the views,
	// this will load the templates.
	//app.RegisterView(tmpl)

	return &HttpServer{
		httpPort:       opts.HttpPort,
		httpHost:       opts.HttpHost,
		appId:          opts.AppId,
		logLevel:       opts.LogLevel,
		daprDddClient:  daprDddClient,
		actorFactories: opts.ActorFactories,
		subscribes:     opts.Subscribes,
		controllers:    opts.Controllers,
		eventTypes:     opts.EventTypes,
		authToken:      opts.AuthToken,
		webRootPath:    opts.WebRootPath,
		envConfig:      opts.EnvConfig,
		app:            app,
		inits:          opts.Inits,
	}

}

func (s *HttpServer) EnvConfig() *EnvConfig {
	return s.envConfig
}

func (s *HttpServer) App() *iris.Application {
	return s.app
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
	for _, opt := range s.inits {
		if err := opt(s); err != nil {
			return err
		}
	}

	/*
		if s.envConfig.App.RsServer.Enable {
			fsManger := s.envConfig.GetFsManager()
			fileFs, fileOk := fsManger.Get(s.envConfig.App.RsServer.FileFsId)
			if !fileOk {
				return errors.New("rsServer file fs key not found")
			}
			httpFs, httpOk := fsManger.Get(s.envConfig.App.RsServer.HttpFsId)
			if !httpOk {
				return errors.New("rsServer http fs key not found")
			}
			fsCfg := rs_server.NewFsConfig(fileFs, httpFs)
			server := rs_server.New(app, fsCfg, true)
			if err := server.Run(); err != nil {
				return err
			}
		}
	*/

	if err := s.addRenderHandler(app); err != nil {
		return err
	}

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
		fmt.Printf("---------- %s running ----------\r\n", s.envConfig.App.AppId)
	}); err != nil {
		return err
	}

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
	fs, ok := s.envConfig.fsManager.Get(fsKey)
	if !ok {
		return errors.New("template fsManager fs key not found %s", fsKey)
	}
	render, err := template.NewHandler(fs)
	if err != nil {
		return err
	}
	apiUrl := fmt.Sprintf("%s/{filePath:path}", s.envConfig.App.Template.ApiUrl)
	app.Get(apiUrl, func(ictx iris.Context) {
		filePath := ictx.URLParam("filePath")
		render(ictx, filePath)
	})
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
