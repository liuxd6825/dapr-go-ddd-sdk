package dbevent

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

func ApplyEvent(ctx context.Context, agg *Aggregate, event *Event, opts ...*ddd.ApplyEventOptions) *dapr.ApplyEventResponse {
	data, err := ddd.ApplyEvent(ctx, agg, event, NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func ApplyEvents(ctx context.Context, agg *Aggregate, events []*Event, opts ...*ddd.ApplyEventOptions) *dapr.ApplyEventResponse {
	data, err := ddd.ApplyEvents(ctx, agg, NewEvents(events), NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func CreateEvent(ctx context.Context, agg *Aggregate, event *Event, opts ...*ddd.ApplyEventOptions) *dapr.CreateEventResponse {
	data, err := ddd.CreateEvent(ctx, agg, event, NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func CreateEvents(ctx context.Context, agg *Aggregate, events []*Event, opts ...*ddd.ApplyEventOptions) *dapr.CreateEventResponse {
	data, err := ddd.CreateEvents(ctx, agg, NewEvents(events), NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func DeleteEvent(ctx context.Context, agg *Aggregate, event *Event, opts ...*ddd.ApplyEventOptions) *dapr.DeleteEventResponse {
	data, err := ddd.DeleteEvent(ctx, agg, event, NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func DeleteEvents(ctx context.Context, agg *Aggregate, events []*Event, opts ...*ddd.ApplyEventOptions) *dapr.DeleteEventResponse {
	data, err := ddd.DeleteEvents(ctx, agg, NewEvents(events), NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func NewEvents(events []*Event) []ddd.DomainEvent {
	var list []ddd.DomainEvent
	for _, e := range events {
		list = append(list, e)
	}
	return list
}

func NewOptions(opts ...*ddd.ApplyEventOptions) *ddd.ApplyEventOptions {
	var closeEventSource = true
	res := ddd.ApplyEventOptions{
		CloseEventSource: &closeEventSource,
	}
	res.Merge(opts...)
	return &res
}

func NewDomainEvent() any {
	return map[string]any{}
}

// GetEventType
//
//	@Description: 获取事件类型
//	@param appId 应用ID
//	@param aggName 聚合根类型名称
//	@param operateType 操作类型 增加，更新，删除等
//	@param eventVersion 版本号
//	@return string
func GetEventType(appId string, aggName string, operateType string) string {
	return fmt.Sprintf("%s.%s.%s", appId, aggName, operateType)
}
