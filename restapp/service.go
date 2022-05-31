package restapp

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/go-sdk/actor"
	"github.com/liuxd6825/go-sdk/actor/config"
	actorErr "github.com/liuxd6825/go-sdk/actor/error"
	"github.com/liuxd6825/go-sdk/actor/runtime"
	"github.com/liuxd6825/go-sdk/service/common"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"net/http"
)

type ServiceOptions struct {
	AppId          string
	HttpHost       string
	HttpPort       int
	LogLevel       applog.Level
	EventStorages  map[string]ddd.EventStorage
	ActorFactories *[]actor.Factory
	Subscribes     *[]RegisterSubscribe
	Controllers    *[]Controller
	EventTypes     *[]RegisterEventType
	AuthToken      string
	WebRootPath    string
}
type service struct {
	app            *iris.Application
	appId          string
	httpHost       string
	httpPort       int
	logLevel       applog.Level
	daprDddClient  daprclient.DaprDddClient
	eventStorages  map[string]ddd.EventStorage
	actorFactories *[]actor.Factory
	subscribes     *[]RegisterSubscribe
	controllers    *[]Controller
	eventTypes     *[]RegisterEventType
	authToken      string
	webRootPath    string
}

func (s *service) AddServiceInvocationHandler(name string, fn common.ServiceInvocationHandler) error {
	panic("implement me")
}

func (s *service) AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error {
	panic("implement me")
}

func (s *service) AddBindingInvocationHandler(name string, fn common.BindingInvocationHandler) error {
	panic("implement me")
}

func (s *service) RegisterActorImplFactory(f actor.Factory, opts ...config.Option) {
	runtime.GetActorRuntimeInstance().RegisterActorFactory(f, opts...)
}

func (s *service) Stop() error {
	panic("implement me")
}

func (s *service) GracefulStop() error {
	panic("implement me")
}

func (s *service) setOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	w.Header().Set("Allow", "POST,OPTIONS")
}

func NewService(daprDddClient daprclient.DaprDddClient, opts *ServiceOptions) common.Service {
	return &service{
		httpPort:       opts.HttpPort,
		httpHost:       opts.HttpHost,
		appId:          opts.AppId,
		logLevel:       opts.LogLevel,
		daprDddClient:  daprDddClient,
		eventStorages:  opts.EventStorages,
		actorFactories: opts.ActorFactories,
		subscribes:     opts.Subscribes,
		controllers:    opts.Controllers,
		eventTypes:     opts.EventTypes,
		authToken:      opts.AuthToken,
		webRootPath:    opts.WebRootPath,
		app:            iris.New(),
	}
}

func (s *service) Start() error {
	app := s.app

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

	// 注册消息订阅
	if s.subscribes != nil {
		for _, subscribe := range *s.subscribes {
			if subscribe != nil {
				if _, err := s.registerSubscribeHandler(subscribe.GetSubscribes(), subscribe.GetHandler()); err != nil {
					return err
				}
			}
		}
	}

	// 注册控制器
	if s.controllers != nil {
		for _, c := range *s.controllers {
			if c != nil {
				s.registerController(s.webRootPath, c)
			}
		}
	}

	// 注册领域事件类型
	if s.eventTypes != nil {
		for _, t := range *s.eventTypes {
			if err := ddd.RegisterEventType(t.EventType, t.Version, t.NewFunc); err != nil {
				return errors.New(fmt.Sprintf("RegisterEventType() error:\"%s\" , EventType=\"%s\", Version=\"%s\"", err.Error(), t.EventType, t.Version))
			}
		}
	}

	// 注册事件存储器
	if s.eventStorages != nil {
		for key, es := range s.eventStorages {
			ddd.RegisterEventStorage(key, es)
		}
	}
	if err := ddd.StartSubscribeHandlers(); err != nil {
		return err
	}

	if s.actorFactories != nil {
		for _, f := range *s.actorFactories {
			s.RegisterActorImplFactory(f)
		}
	}

	if err := app.Run(iris.Addr(fmt.Sprintf("%s:%d", s.httpHost, s.httpPort))); err != nil {
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

func (s *service) subscribesHandler(ctx *context.Context) {
	data := ddd.GetSubscribes()
	_, _ = ctx.JSON(data)
}

func (s *service) eventTypesHandler(context *context.Context) {

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

//
// registerQueryHandler
// @Description: 注册领域事件控制器
// @param handlers
// @return error
//
func (s *service) registerQueryHandler(handlers ...ddd.SubscribeHandler) error {
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
func (s *service) registerSubscribeHandler(subscribes *[]ddd.Subscribe, queryEventHandler ddd.QueryEventHandler) (ddd.SubscribeHandler, error) {
	handler := ddd.NewSubscribeHandler(subscribes, queryEventHandler, func(sh ddd.SubscribeHandler, subscribe ddd.Subscribe) (err error) {
		defer func() {
			if e := ddd_errors.GetRecoverError(recover()); e != nil {
				err = e
			}
		}()
		s.app.Handle("POST", subscribe.Route, func(c *context.Context) {
			println(fmt.Sprintf("topic:%s; route:%s;", subscribe.Topic, subscribe.Route))
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

//
// RegisterRestController
// @Description: 注册UserInterface层Controller
// @param relativePath
// @param configurators
//
func (s *service) registerController(relativePath string, controllers ...Controller) {
	if controllers == nil && len(controllers) == 0 {
		return
	}
	configurators := func(app *mvc.Application) {
		for _, c := range controllers {
			app.Handle(c)
		}
	}
	mvc.Configure(s.app.Party(relativePath), configurators)
}
