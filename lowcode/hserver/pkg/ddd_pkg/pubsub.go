package ddd_pkg

import (
	"context"
	"github.com/dapr/go-sdk/client"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
)

type Pubsub struct {
	server element.Server
	client dapr.DaprClient
}

func New(server element.Server) *Pubsub {
	cli := server.HttpServer().DaprClient()
	return &Pubsub{
		server: server,
		client: cli,
	}
}

func (p *Pubsub) ApplyEvent(ctx context.Context, agg *Aggregate, event *DomainEvent, opts ...*ApplyEventOptions) (*dapr.ApplyEventResponse, error) {
	return ddd.ApplyEvent(ctx, agg, event, opts...)
}

func (p *Pubsub) ApplyEvents(ctx context.Context, agg *Aggregate, events []*DomainEvent, opts ...*ApplyEventOptions) (*dapr.ApplyEventResponse, error) {
	return ddd.ApplyEvents(ctx, agg, newEvents(events), opts...)
}

func (p *Pubsub) CreateEvent(ctx context.Context, agg *Aggregate, event *DomainEvent, opts ...*ApplyEventOptions) (*dapr.CreateEventResponse, error) {
	return ddd.CreateEvent(ctx, agg, event, opts...)
}

func (p *Pubsub) CreateEvents(ctx context.Context, agg *Aggregate, events []*DomainEvent, opts ...*ApplyEventOptions) (*dapr.CreateEventResponse, error) {
	return ddd.CreateEvents(ctx, agg, newEvents(events), opts...)
}

func (p *Pubsub) DeleteEvent(ctx context.Context, agg *Aggregate, event *DomainEvent, opts ...*ApplyEventOptions) (*dapr.DeleteEventResponse, error) {
	return ddd.DeleteEvent(ctx, agg, event, opts...)
}

func (p *Pubsub) DeleteEvents(ctx context.Context, agg *Aggregate, events []*DomainEvent, opts ...*ApplyEventOptions) (*dapr.DeleteEventResponse, error) {
	return ddd.DeleteEvents(ctx, agg, newEvents(events), opts...)
}

// Subscribe 订阅事件
func (p *Pubsub) Subscribe(ctx context.Context, opt client.SubscriptionOptions, handler client.SubscriptionHandleFunction) any {
	fun, err := p.client.SubscribeWithHandler(ctx, opt, handler)
	if err != nil {
		panic(err)
	}
	return fun
}

func newEvents(events []*DomainEvent) []ddd.DomainEvent {
	list := make([]ddd.DomainEvent, len(events))
	for i, e := range events {
		list[i] = e
	}
	return list
}
