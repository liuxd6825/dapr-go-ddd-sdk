package event

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"time"
)

type EventData struct {
	EventId     string         `json:"eventId" validate:"required"`   // 领域事件ID
	CommandId   string         `json:"commandId" validate:"required"` // 关联命令ID
	CreatedTime time.Time      `json:"time" validate:"required"`      // 事件创建时间
	Version     string         `json:"version" validate:"required"`   // 事件版本
	EventType   string         `json:"eventType" validate:"required"` // 事件类型
	Data        map[string]any `json:"data" validate:"required"`      // 业务字段项
}

type Handle struct {
	service element.Service
	fun     element.Func
	tag     *element.FuncTag
	config  *Config
}

func Register(service element.Service, fun element.Func, tag *element.FuncTag) {
	h, err := NewHandle(service, fun, tag)
	if err != nil {
		panic(err)
	}
	if err = h.Initialize(); err != nil {
		panic(err)
	}
}

func NewHandle(service element.Service, fun element.Func, tag *element.FuncTag) (*Handle, error) {
	h := &Handle{
		service: service,
		fun:     fun,
		tag:     tag,
		config:  NewConfig(tag),
	}
	return h, nil
}

func (e *Handle) Close() error {
	return nil
}

func (e *Handle) newEvent() interface{} {
	return &EventData{}
}

func (e *Handle) CallEventHandler(ctx context.Context, handler any, eventType string, eventVersion string, event any, metadata ddd.Metadata) error {
	return e.run(ctx, eventType, eventVersion, event, metadata)
}

func (e *Handle) run(ctx context.Context, eventType string, eventVersion string, event any, metadata ddd.Metadata) error {
	values := &element.ApiRunValues{
		Server:   e.service.Server(),
		Self:     e.service,
		WorkPath: e.WorkPath(),
	}
	tenantId := metadata["tenantId"]
	_, err := e.service.Funcs().Run(e.fun.Config().FuncName, values, true, func(vm *goja.Runtime) error {
		_ = vm.Set("tenantId", tenantId)
		_ = vm.Set("ctx", ctx)
		_ = vm.Set("event", event)
		_ = vm.Set("eventVersion", eventVersion)
		_ = vm.Set("eventType", eventType)
		_ = vm.Set("metadata", metadata)
		return nil
	})
	return err
}

func (e *Handle) WorkPath() string {
	return e.service.WorkPath()
}

func (e *Handle) Initialize() error {
	version := e.config.Ver()
	appId := e.config.AppId()
	pubsub := e.config.Pubsub()
	eventType := e.config.Name()

	if err := ddd.RegisterEventType(eventType, version, e.newEvent); err != nil {
		return errors.New(fmt.Sprintf("RegisterEventType() error:\"%s\" , EventType=\"%s\", Version=\"%s\"", err.Error(), eventType, version))
	}
	routeUrl := fmt.Sprintf("domain-event/pubsub/%s/app-id/%s/event-type/%s/ver/%s", pubsub, appId, eventType, version)
	sub := &ddd.Subscribe{
		PubsubName: pubsub,
		Topic:      eventType,
		Route:      routeUrl,
	}
	_, err := e.service.Server().HttpServer().RegisterSubscribeHandler(
		[]*ddd.Subscribe{sub},
		e,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
