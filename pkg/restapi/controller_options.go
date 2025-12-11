package restapi

import (
	context2 "context"
	"reflect"

	"github.com/kataras/iris/v12/context"
)

type GetFieldValueFunc func(ictx *context.Context, parentObject any, fieldType reflect.StructField, fieldValue reflect.Value) (outFieldValue any, ok bool, err error)
type BeforeFunc func(ictx *context.Context, params any) (any, context2.Context, error)
type AfterFunc func(ctx context2.Context, ictx *context.Context, data any, err error) (any, error)
type APIOptions struct {
	Description      string                       // 说明
	IsAuthentication *bool                        // 是否进行身份认证
	ParamsInBody     *bool                        // 强制要求参数入HttpBody中获取
	TranDBKeys       []string                     // 数据库事物
	InitMethod       func(callMethod *CallMethod) // 初始化函数
	GetFieldValue    GetFieldValueFunc
	Before           BeforeFunc
	After            AfterFunc
	Event            *EventOption
	ApiService       any // API服务
}

type EventOption struct {
	Pubsub string    `json:"pubsub"`
	Topic  string    `json:"topic"`
	Meta   EventMeta `json:"meta"`
}

// dapr消息
type EventMeta struct {
	// <--非dapr消息，关系配置 "true"
	RawPayload bool
	// 告诉 Dapr 直接消费这个名字的队列 "record-import-master-event"
	QueueName string
	// 明确告诉 Dapr 在绑定队列时使用哪个路由键
	// 这个值必须与消息发布方使用的路由键相匹配
	//  "record-import-master-event"
	RoutingKey string
}

func WithEvent(pubsub string, topic string) APIOptions {
	event := &EventOption{
		Pubsub: pubsub,
		Topic:  topic,
	}
	return APIOptions{
		Event: event,
	}
}

func WithEventMeta(pubsub string, topic string, meta *EventMeta) APIOptions {
	event := &EventOption{
		Pubsub: pubsub,
		Topic:  topic,
	}
	if meta != nil {
		event.Meta = *meta
	}
	return APIOptions{
		Event: event,
	}
}

func NewEventMeta() *EventMeta {
	return &EventMeta{}
}

func (m *EventMeta) SetRawPayload(value bool) *EventMeta {
	m.RawPayload = value
	return m
}

func (m *EventMeta) SetQueueName(value string) *EventMeta {
	m.QueueName = value
	return m
}

func (m *EventMeta) SetRoutingKey(value string) *EventMeta {
	m.RoutingKey = value
	return m
}

func WithIsAuthentication(val bool) APIOptions {
	return APIOptions{
		IsAuthentication: &val,
	}
}

func WithParamsInBody(val bool) APIOptions {
	return APIOptions{
		ParamsInBody: &val,
	}
}

func WithInitMethod(method func(callMethod *CallMethod)) APIOptions {
	return APIOptions{
		InitMethod: method,
	}
}

func WithGetFieldValue(method GetFieldValueFunc) APIOptions {
	return APIOptions{
		GetFieldValue: method,
	}
}

func WithBefore(method BeforeFunc) APIOptions {
	return APIOptions{
		Before: method,
	}
}

func WithAfter(method AfterFunc) APIOptions {
	return APIOptions{
		After: method,
	}
}

func WithTranDbKey(dbKeys ...string) APIOptions {
	return APIOptions{
		TranDBKeys: dbKeys,
	}
}

func NewAPIOptions(opts ...APIOptions) APIOptions {
	o := APIOptions{}
	for _, i := range opts {
		if i.InitMethod != nil {
			o.InitMethod = i.InitMethod
		}
		if i.GetFieldValue != nil {
			o.GetFieldValue = i.GetFieldValue
		}
		if i.Before != nil {
			o.Before = i.Before
		}
		if i.After != nil {
			o.After = i.After
		}
		if i.ParamsInBody != nil {
			o.ParamsInBody = i.ParamsInBody
		}
		if i.IsAuthentication != nil {
			o.IsAuthentication = i.IsAuthentication
		}
		if i.Event != nil {
			o.Event = i.Event
		}
		if i.Description != "" {
			o.Description = i.Description
		}
		if len(i.TranDBKeys) > 0 {
			o.TranDBKeys = i.TranDBKeys
		}
	}
	return o
}

func (o *APIOptions) GetIsAuthentication() bool {
	if o.IsAuthentication == nil {
		return false
	}
	return *o.IsAuthentication
}

func (o *APIOptions) SetTranDBKeys(dbKeys ...string) *APIOptions {
	o.TranDBKeys = dbKeys
	return o
}

func (o *APIOptions) GetTranDBKeys() []string {
	return o.TranDBKeys
}

func (o *APIOptions) SetEvent(pubsub string, topic string) *APIOptions {
	if pubsub == "" || topic == "" {
		panic("pubsub or topic is empty")
	}
	o.Event = &EventOption{
		Pubsub: pubsub,
		Topic:  topic,
	}
	return o
}

func (o *APIOptions) GetEvent() *EventOption {
	return o.Event
}

func (o *APIOptions) SetDescription(description string) *APIOptions {
	o.Description = description
	return o
}

func (o *APIOptions) GetDescription() string {
	return o.Description
}
