package sub

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
)

type SubEvent struct {
	element.Base
	config      *element.SubEventConfig
	server      element.Server
	service     element.SubService
	srcFileName string
}

func NewSubEvent(server element.Server, service element.SubService, srcFileName string, config *element.SubEventConfig) (element.SubEvent, error) {
	var err error
	r := &SubEvent{
		server:      server,
		service:     service,
		config:      config,
		srcFileName: srcFileName,
	}

	fsOpts := fsopts.NewOptionsWidthFileName(srcFileName, server.RootPath())
	r.Base, err = server.Factory().NewBase(srcFileName, server.Logger(), server, server.Factory(), fsOpts)
	if err != nil {
		return nil, err
	}
	r.SetPkg(server.Pkg())
	return r, nil

}

func (e *SubEvent) Config() *element.SubEventConfig {
	return e.config
}

func (e *SubEvent) Close() error {
	return nil
}

func (e *SubEvent) newEvent() interface{} {
	return map[string]any{}
}

func (e *SubEvent) CallEventHandler(ctx context.Context, handler any, eventType string, eventVersion string, event any, metadata ddd.Metadata) error {
	return e.runScript(ctx, eventType, eventVersion, event, metadata)
}

func (e *SubEvent) runScript(ctx context.Context, eventType string, eventVersion string, event any, metadata ddd.Metadata) error {
	values := &element.ApiRunValues{
		Server:   e.server,
		Self:     e.service,
		WorkPath: e.WorkPath(),
	}
	tenantId := metadata["tenantId"]
	_, err := e.Scripts().RunScript(e.config.Script.FuncName, values, true, func(vm *goja.Runtime) error {
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

func (e *SubEvent) WorkPath() string {
	return ""
}

func (e *SubEvent) Initialize() error {
	eventType := e.config.Name
	version := e.config.Version
	appId := e.service.Config().AppId
	pubsub := e.service.Config().Pubsub

	if err := ddd.RegisterEventType(e.config.Name, e.config.Version, e.newEvent); err != nil {
		return errors.New(fmt.Sprintf("RegisterEventType() error:\"%s\" , EventType=\"%s\", Version=\"%s\"", err.Error(), eventType, version))
	}
	routeUrl := fmt.Sprintf("domain-event/%s/%s/%s/%s", pubsub, appId, e.config.Name, e.config.Version)
	sub := &ddd.Subscribe{
		PubsubName: pubsub,
		Topic:      e.config.Name,
		Route:      routeUrl,
	}
	_, err := e.server.HttpServer().RegisterSubscribeHandler(
		[]*ddd.Subscribe{sub},
		e,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
