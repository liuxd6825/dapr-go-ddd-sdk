package restapi

import (
	context2 "context"
	"reflect"

	"github.com/kataras/iris/v12/context"
)

type GetFieldValueFunc func(ictx *context.Context, parentObject any, fieldType reflect.StructField, fieldValue reflect.Value) (outFieldValue any, ok bool, err error)
type BeforeFunc func(ictx *context.Context, params any) (any, context2.Context, error)
type AfterFunc func(ctx context2.Context, ictx *context.Context, data any, err error) (any, error)
type CallOptions struct {
	IsAuthentication *bool                        // 是否进行身份认证
	ParamsInBody     *bool                        // 强制要求参数入HttpBody中获取
	InitMethod       func(callMethod *CallMethod) // 初始化函数
	GetFieldValue    GetFieldValueFunc
	Before           BeforeFunc
	After            AfterFunc
	Event            *EventOption
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

func WithEvent(pubsub string, topic string) CallOptions {
	event := &EventOption{
		Pubsub: pubsub,
		Topic:  topic,
	}
	return CallOptions{
		Event: event,
	}
}

func WithEventMeta(pubsub string, topic string, meta *EventMeta) CallOptions {
	event := &EventOption{
		Pubsub: pubsub,
		Topic:  topic,
	}
	if meta != nil {
		event.Meta = *meta
	}
	return CallOptions{
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

func WithIsAuthentication(val bool) CallOptions {
	return CallOptions{
		IsAuthentication: &val,
	}
}

func WithParamsInBody(val bool) CallOptions {
	return CallOptions{
		ParamsInBody: &val,
	}
}

func WithInitMethod(method func(callMethod *CallMethod)) CallOptions {
	return CallOptions{
		InitMethod: method,
	}
}

func WithGetFieldValue(method GetFieldValueFunc) CallOptions {
	return CallOptions{
		GetFieldValue: method,
	}
}

func WithBefore(method BeforeFunc) CallOptions {
	return CallOptions{
		Before: method,
	}
}

func WithAfter(method AfterFunc) CallOptions {
	return CallOptions{
		After: method,
	}
}

func NewCallOptions(opts ...CallOptions) CallOptions {
	o := CallOptions{}
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
	}
	return o
}

func (o *CallOptions) GetIsAuthentication() bool {
	if o.IsAuthentication == nil {
		return false
	}
	return *o.IsAuthentication
}

func (o *CallOptions) SetEvent(pubsub string, topic string) *CallOptions {
	if pubsub == "" || topic == "" {
		panic("pubsub or topic is empty")
	}
	o.Event = &EventOption{
		Pubsub: pubsub,
		Topic:  topic,
	}
	return o
}

func (o *CallOptions) GetEvent() *EventOption {
	return o.Event
}
