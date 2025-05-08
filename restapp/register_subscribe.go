package restapp

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type RegisterSubscribeOptions struct {
	interceptors []ddd.SubscribeInterceptorFunc
}

type RegisterSubscribe interface {
	GetSubscribes() []*ddd.Subscribe
	GetHandler() ddd.QueryEventHandler
	GetInterceptor() []ddd.SubscribeInterceptorFunc
}

type registerSubscribe struct {
	subscribes   []*ddd.Subscribe
	handler      ddd.QueryEventHandler
	interceptors []ddd.SubscribeInterceptorFunc
}

var _subscribeInterceptor []ddd.SubscribeInterceptorFunc

func RegisterSubscribeInterceptor(items ...ddd.SubscribeInterceptorFunc) {
	_subscribeInterceptor = append(_subscribeInterceptor, items...)
}

func NewRegisterSubscribeOptions(opts ...*RegisterSubscribeOptions) *RegisterSubscribeOptions {
	o := &RegisterSubscribeOptions{}
	for _, item := range opts {
		if item.interceptors != nil {
			o.interceptors = item.interceptors
		}
	}
	return o
}

func NewRegisterSubscribe(subscribes []*ddd.Subscribe, handler ddd.QueryEventHandler, options ...*RegisterSubscribeOptions) RegisterSubscribe {
	opt := NewRegisterSubscribeOptions(options...)
	return &registerSubscribe{
		subscribes:   subscribes,
		handler:      handler,
		interceptors: opt.interceptors,
	}
}

func (r *registerSubscribe) GetInterceptor() []ddd.SubscribeInterceptorFunc {
	return r.interceptors
}

func (r *registerSubscribe) GetSubscribes() []*ddd.Subscribe {
	return r.subscribes
}

func (r *registerSubscribe) GetHandler() ddd.QueryEventHandler {
	return r.handler
}

func (o *RegisterSubscribeOptions) SetInterceptors(v []ddd.SubscribeInterceptorFunc) *RegisterSubscribeOptions {
	o.interceptors = v
	return o
}
