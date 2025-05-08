package server

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

// EventHandler 事件处理器
type EventHandler struct {
	VM         *goja.Runtime
	Object     *goja.Object
	subscribes []*Subscribe // 处理器所订阅的消息列表
}

// CallEventHandler
//
//	@Description: 执行消息处理函数
//	@receiver h
//	@param ctx
//	@param handler 事件处理器
//	@param eventType    事件类型
//	@param eventVersion 事件版本
//	@param event 事件数据
//	@param metadata 事件元数据
//	@return err
func (h *EventHandler) CallEventHandler(ctx context.Context, handler any, eventType string, eventVersion string, event any, metadata ddd.Metadata) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
	}()
	methodName := h.GetMethodName(eventType)
	v := h.Object.Get(methodName)
	if goja.IsUndefined(v) || goja.IsNaN(v) || goja.IsNull(v) {
		return errors.New("method: %s is not defined", methodName)
	}
	fun, ok := goja.AssertFunction(v)
	if !ok {
		return errors.New("rs-server subscribe handler not fount %s()", methodName)
	}
	c := h.VM.ToValue(ctx)
	e := h.VM.ToValue(event)
	m := h.VM.ToValue(metadata)
	_, vErr := fun(h.Object, c, e, m)
	if vErr != nil {
		return vErr
	}
	return nil
}

func (h *EventHandler) GetMethodName(eventType string) string {
	for _, s := range h.subscribes {
		if s.Topic == eventType {
			return s.FuncName
		}
	}
	return ""
}

type Subscribe = ddd.Subscribe

type SubscribeItem struct {
	AppId        string            `json:"appId"`
	PubsubName   string            `json:"pubsubName"`
	EventType    string            `json:"eventType"`
	EventVersion string            `json:"eventVersion"`
	Route        string            `json:"route,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	FuncName     string            `json:"funcName"`
}

type RegisterSubscribe = restapp.RegisterSubscribe

type RegisterSubscribeOptions struct {
	interceptors []ddd.SubscribeInterceptorFunc
}

func newRegisterSubscribeOptions(options ...*RegisterSubscribeOptions) *RegisterSubscribeOptions {
	r := &RegisterSubscribeOptions{}
	for _, option := range options {
		if option != nil && option.interceptors != nil {
			r.interceptors = append(r.interceptors, option.interceptors...)
		}
	}
	return r
}

type subscribeService struct {
	subscribes   []*Subscribe
	handler      ddd.QueryEventHandler
	interceptors []ddd.SubscribeInterceptorFunc
}

var (
	registerSubscribes []RegisterSubscribe
)

func RegisterSubscribeService(appId string, subscribeItems []*SubscribeItem, vm *goja.Runtime, serviceObj *goja.Object, options ...*RegisterSubscribeOptions) RegisterSubscribe {
	subscribes := newSubscribes(appId, subscribeItems)
	service := &subscribeService{
		subscribes:   subscribes,
		handler:      &EventHandler{Object: serviceObj, VM: vm, subscribes: subscribes},
		interceptors: newRegisterSubscribeOptions(options...).interceptors,
	}
	registerSubscribes = append(registerSubscribes, service)
	return service
}

func newSubscribes(appId string, subscribeItems []*SubscribeItem) []*Subscribe {
	var subscribes []*Subscribe
	for _, sub := range subscribeItems {
		aId := sub.AppId
		if aId == "" {
			aId = appId
		}
		subscribe := &Subscribe{
			PubsubName: sub.PubsubName,
			Topic:      fmt.Sprintf("%s.%s", aId, sub.EventType),
			Route:      sub.Route,
			Metadata:   sub.Metadata,
			FuncName:   sub.FuncName,
		}
		subscribes = append(subscribes, subscribe)
	}
	return subscribes
}

func GetRegisterSubscribe() []RegisterSubscribe {
	return registerSubscribes
}

func (r *subscribeService) GetSubscribes() []*ddd.Subscribe {
	return r.subscribes
}

func (r *subscribeService) GetHandler() ddd.QueryEventHandler {
	return r.handler
}

func (r *subscribeService) GetInterceptor() []ddd.SubscribeInterceptorFunc {
	return r.interceptors
}
