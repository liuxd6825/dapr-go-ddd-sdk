package ddd

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
)

// Subscribe internally represents single topic subscription.
type Subscribe struct {
	// PubsubName is name of the pub/sub this message came from.
	PubsubName string `json:"pubsubname"`
	// Topic is the name of the topic.
	Topic string `json:"topic"`
	// Route is the route of the handler where HTTP topic events should be published (passed as Path in gRPC).
	Route string `json:"route,omitempty"`
	// Routes specify multiple routes where topic events should be sent.
	Routes *TopicRoutes `json:"routes,omitempty"`
	// Metadata is the subscription metadata.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// TopicRoutes encapsulates the default route and multiple routing rules.
type TopicRoutes struct {
	Rules   []TopicRule `json:"rules,omitempty"`
	Default string      `json:"default,omitempty"`

	// priority is used to track duplicate priorities where priority > 0.
	// when priority is not specified (0), then the order in which they
	// were added is used.
	priorities map[int]struct{}
}

// TopicRule represents a single routing rule.
type TopicRule struct {
	// Match is the CEL expression to match on the CloudEvent envelope.
	Match string `json:"match"`
	// Path is the HTTP path to post the event to (passed as Path in gRPC).
	Path string `json:"path"`
	// priority is the optional priority order (low to high) for this rule.
	priority int `json:"-"`
}

// NewSubscribe 新建消息订阅项
func NewSubscribe(pubsubName string, topic string, route string, metadata map[string]string, handler interface{}) *Subscribe {
	return &Subscribe{
		PubsubName: pubsubName,
		Topic:      topic,
		Route:      route,
		Metadata:   metadata,
	}
}

// SubscribeHandler 消息订阅处理器
type SubscribeHandler interface {
	GetSubscribes() ([]*Subscribe, error)
	RegisterSubscribe(subscribe *Subscribe) error
	SubscribeHandler(ctx context.Context, sCtx SubscribeContext) error
}

type SubscribeContext interface {
	GetBody() ([]byte, error)
	SetErr(err error)
}

type SubscribeInterceptor interface {
	Interceptor(ctx context.Context, sCtx SubscribeContext) (bool, error)
}

type SubscribeHandlerFunc func(sh SubscribeHandler, subscribe *Subscribe) error
type SubscribeInterceptorFunc func(ctx context.Context, sctx SubscribeContext) (bool, error)

// SubscribeHandler 消息订阅处理器
type subscribeHandler struct {
	subscribes           []*Subscribe
	queryEventHandler    QueryEventHandler
	subscribeHandlerFunc SubscribeHandlerFunc
	interceptors         []SubscribeInterceptorFunc
}

func NewSubscribeHandler(subscribes []*Subscribe, queryEventHandler QueryEventHandler, subscribeHandlerFunc SubscribeHandlerFunc, interceptors []SubscribeInterceptorFunc) SubscribeHandler {
	return &subscribeHandler{
		subscribes:           subscribes,
		queryEventHandler:    queryEventHandler,
		subscribeHandlerFunc: subscribeHandlerFunc,
		interceptors:         interceptors,
	}
}

func (h *subscribeHandler) GetSubscribes() ([]*Subscribe, error) {
	return h.subscribes, nil
}

func (h *subscribeHandler) RegisterSubscribe(subscribe *Subscribe) error {
	return h.subscribeHandlerFunc(h, subscribe)
}

// SubscribeHandler
//
//	@Description: 消息订阅处理器
//	@receiver h
//	@param ctx
//	@param sctx
//	@return error
func (h *subscribeHandler) SubscribeHandler(ctx context.Context, sctx SubscribeContext) error {
	cancel, err := h.interceptor(ctx, sctx)
	if cancel || err != nil {
		return err
	}
	data, err := sctx.GetBody()
	if err != nil {
		return err
	}
	return daprclient.NewEventRecordByJsonBytes(data).OnSuccess(ctx, func(ctx context.Context, r *daprclient.EventRecord) error {
		return CallEventHandler(ctx, h.queryEventHandler, r)
	}).OnError(func(err error, r *daprclient.EventRecord) {
		if r != nil {
			fields := logs.Fields{
				"logFunc":   "SubscribeHandler",
				"eventId":   r.EventId,
				"eventType": r.EventType,
				"event":     r.EventVersion,
				"eventData": r.EventData,
				"error":     err.Error(),
			}
			logs.Error(ctx, r.TenantId, fields)
		} else {
			logs.Errorf(ctx, "", nil, "领域事件处理错误: data:%s; error:%s", string(data), err.Error())
		}

	}).GetError()
}

// interceptor
//
//	@Description: 消息拦截器
//	@receiver h
//	@param ctx 上下文
//	@param sctx dapr消息上下文
//	@return bool true:已拦截不需要后续处理。
//	@return error 错误
func (h *subscribeHandler) interceptor(ctx context.Context, sctx SubscribeContext) (cancel bool, err error) {
	defer func() {
		if e := recover(); e != nil {
			if err := e.(error); err != nil {
				return
			}
		}
	}()
	err = nil
	cancel = false
	interceptor, ok := h.queryEventHandler.(SubscribeInterceptor)
	if ok {
		cancel, err = interceptor.Interceptor(ctx, sctx)
	}
	if h.interceptors != nil {
		for _, item := range h.interceptors {
			if c, e := item(ctx, sctx); c {
				cancel = true
			} else if e != nil {

			}
		}
	}
	return
}
