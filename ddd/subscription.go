package ddd

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

// Subscribe dapr消息订阅项
type Subscribe struct {
	PubsubName string            `json:"pubsubName"`
	Topic      string            `json:"topic"`
	Route      string            `json:"route"`
	Metadata   map[string]string `json:"metadata"`
}

// NewSubscribe 新建消息订阅项
func NewSubscribe(pubsubName string, topic, route string, metadata map[string]string, handler interface{}) *Subscribe {
	return &Subscribe{
		PubsubName: pubsubName,
		Topic:      topic,
		Metadata:   metadata,
		Route:      route,
	}
}

// SubscribeHandler 消息订阅处理器
type SubscribeHandler interface {
	GetSubscribes() (*[]Subscribe, error)
	RegisterSubscribe(subscribe Subscribe) error
	CallQueryEventHandler(ctx context.Context, sctx SubscribeContext) error
}

type SubscribeContext interface {
	GetBody() ([]byte, error)
	SetErr(err error)
}

type SubscribeHandlerFunc func(sh SubscribeHandler, subscribe Subscribe) error

// SubscribeHandler 消息订阅处理器
type subscribeHandler struct {
	subscribes           *[]Subscribe
	queryEventHandler    QueryEventHandler
	subscribeHandlerFunc SubscribeHandlerFunc
}

func NewSubscribeHandler(subscribes *[]Subscribe, queryEventHandler QueryEventHandler, subscribeHandlerFunc SubscribeHandlerFunc) SubscribeHandler {
	return &subscribeHandler{
		subscribes:           subscribes,
		queryEventHandler:    queryEventHandler,
		subscribeHandlerFunc: subscribeHandlerFunc,
	}
}

func (h *subscribeHandler) GetSubscribes() (*[]Subscribe, error) {
	return h.subscribes, nil
}

func (h *subscribeHandler) RegisterSubscribe(subscribe Subscribe) error {
	return h.subscribeHandlerFunc(h, subscribe)
}

// CallQueryEventHandler 消息订阅处理器
func (h *subscribeHandler) CallQueryEventHandler(ctx context.Context, sctx SubscribeContext) error {
	data, err := sctx.GetBody()
	if err != nil {
		return err
	}
	return daprclient.NewEventRecordByJsonBytes(data).OnSuccess(func(eventRecord *daprclient.EventRecord) error {
		return CallEventHandler(ctx, h.queryEventHandler, eventRecord)
	}).GetError()
}
