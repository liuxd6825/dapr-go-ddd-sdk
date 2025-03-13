package events

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

func ApplyEvent(ctx context.Context, agg *Aggregate, event *common.Event, opts ...*ddd.ApplyEventOptions) *dapr.ApplyEventResponse {
	data, err := ddd.ApplyEvent(ctx, agg, event, NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func ApplyEvents(ctx context.Context, agg *Aggregate, events []*common.Event, opts ...*ddd.ApplyEventOptions) *dapr.ApplyEventResponse {
	data, err := ddd.ApplyEvents(ctx, agg, NewEvents(events), NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func CreateEvent(ctx context.Context, agg *Aggregate, event *common.Event, opts ...*ddd.ApplyEventOptions) *dapr.CreateEventResponse {
	data, err := ddd.CreateEvent(ctx, agg, event, NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func CreateEvents(ctx context.Context, agg *Aggregate, events []*common.Event, opts ...*ddd.ApplyEventOptions) *dapr.CreateEventResponse {
	data, err := ddd.CreateEvents(ctx, agg, NewEvents(events), NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func DeleteEvent(ctx context.Context, agg *Aggregate, event *common.Event, opts ...*ddd.ApplyEventOptions) *dapr.DeleteEventResponse {
	data, err := ddd.DeleteEvent(ctx, agg, event, NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func DeleteEvents(ctx context.Context, agg *Aggregate, events []*common.Event, opts ...*ddd.ApplyEventOptions) *dapr.DeleteEventResponse {
	data, err := ddd.DeleteEvents(ctx, agg, NewEvents(events), NewOptions(opts...))
	if err != nil {
		panic(err)
	}
	return data
}

func GetEventType(aggName string, opType string) string {
	return common.GetEventType(restapp.GetEnvConfig().GetAppId(), aggName, opType)
}

func NewEvents(events []*common.Event) []ddd.DomainEvent {
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
