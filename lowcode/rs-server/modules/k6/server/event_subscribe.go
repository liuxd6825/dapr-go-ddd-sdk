package server

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type Subscribe = ddd.Subscribe

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

func RegisterSubscribeService(subscribes []*Subscribe, vm *goja.Runtime, serviceObj *goja.Object, options ...*RegisterSubscribeOptions) RegisterSubscribe {
	service := &subscribeService{
		subscribes:   subscribes,
		handler:      &EventHandler{Object: serviceObj, VM: vm, subscribes: subscribes},
		interceptors: newRegisterSubscribeOptions(options...).interceptors,
	}
	registerSubscribes = append(registerSubscribes, service)
	return service
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
