package server

import (
	"context"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"sync"
)

type EventPkg struct {
	vm *goja.Runtime
}

var _eventPkg *EventPkg
var _eventPkgOnce sync.Once

func GetEventPkg() *EventPkg {
	return _eventPkg
}

func CreateEventPkg(vm *goja.Runtime) *EventPkg {
	_eventPkgOnce.Do(func() {
		_eventPkg = NewEventPkg(vm)
	})
	if vm != nil {
		_eventPkg.vm = vm
	}
	return _eventPkg
}

type Aggregate = common.Aggregate

func NewEventPkg(vm *goja.Runtime) *EventPkg {
	return &EventPkg{vm: vm}
}

func (e *EventPkg) ApplyEvent(ctx context.Context, agg *Aggregate, event *common.Event, opts ...*ddd.ApplyEventOptions) *common.Result[*dapr.ApplyEventResponse] {
	data, err := ddd.ApplyEvent(ctx, agg, event, e.newOptions(opts...))
	return common.NewResult[*dapr.ApplyEventResponse](data, err)
}

func (e *EventPkg) ApplyEvents(ctx context.Context, agg *Aggregate, events []*common.Event, opts ...*ddd.ApplyEventOptions) *common.Result[*dapr.ApplyEventResponse] {
	data, err := ddd.ApplyEvents(ctx, agg, e.newEvents(events), e.newOptions(opts...))
	return common.NewResult[*dapr.ApplyEventResponse](data, err)
}

func (e *EventPkg) CreateEvent(ctx context.Context, agg *Aggregate, event *common.Event, opts ...*ddd.ApplyEventOptions) *common.Result[*dapr.CreateEventResponse] {
	data, err := ddd.CreateEvent(ctx, agg, event, e.newOptions(opts...))
	return common.NewResult[*dapr.CreateEventResponse](data, err)
}

func (e *EventPkg) CreateEvents(ctx context.Context, agg *Aggregate, events []*common.Event, opts ...*ddd.ApplyEventOptions) *common.Result[*dapr.CreateEventResponse] {
	data, err := ddd.CreateEvents(ctx, agg, e.newEvents(events), e.newOptions(opts...))
	return common.NewResult[*dapr.CreateEventResponse](data, err)
}

func (e *EventPkg) DeleteEvent(ctx context.Context, agg *Aggregate, event *common.Event, opts ...*ddd.ApplyEventOptions) *common.Result[*dapr.DeleteEventResponse] {
	data, err := ddd.DeleteEvent(ctx, agg, event, e.newOptions(opts...))
	return common.NewResult[*dapr.DeleteEventResponse](data, err)
}

func (e *EventPkg) DeleteEvents(ctx context.Context, agg *Aggregate, events []*common.Event, opts ...*ddd.ApplyEventOptions) *common.Result[*dapr.DeleteEventResponse] {
	data, err := ddd.DeleteEvents(ctx, agg, e.newEvents(events), e.newOptions(opts...))
	return common.NewResult[*dapr.DeleteEventResponse](data, err)
}

func (e *EventPkg) AddEventHandler(subscribes []*Subscribe, serviceObj *goja.Object, options ...*RegisterSubscribeOptions) {
	RegisterSubscribeService(subscribes, e.vm, serviceObj, options...)
}

func (e *EventPkg) newEvents(events []*common.Event) []ddd.DomainEvent {
	var list []ddd.DomainEvent
	for _, e := range events {
		list = append(list, e)
	}
	return list
}

func (e *EventPkg) newOptions(opts ...*ddd.ApplyEventOptions) *ddd.ApplyEventOptions {
	var closeEventSource = true
	res := ddd.ApplyEventOptions{
		CloseEventSource: &closeEventSource,
	}
	res.Merge(opts...)
	return &res
}
