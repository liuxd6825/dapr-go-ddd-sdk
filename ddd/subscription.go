package ddd

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
)

// Subscribe dapr消息订阅项
type Subscribe struct {
	PubsubName string            `json:"pubsubName"`
	Topic      string            `json:"topic,omitempty"`
	Routes     *Routes           `json:"routes,omitempty"` // map[string]string
	Route      string            `json:"route,omitempty"`  // map[string]string
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Routes encapsulates the rules and optional default path for a topic.
type Routes struct {
	// The list of rules for this topic.
	// +optional
	Rules []Rule `json:"rules,omitempty"`
	// The default path for this topic.
	// +optional
	Default string `json:"default,omitempty"`
}

// Rule is used to specify the condition for sending
// a message to a specific path.
type Rule struct {
	// The optional CEL expression used to match the event.
	// If the match is not specified, then the route is considered
	// the default. The rules are tested in the order specified,
	// so they should be define from most-to-least specific.
	// The default route should appear last in the list.
	Match string `json:"match"`

	// The path for events that match this rule.
	Path string `json:"path"`
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
