package events_pkg

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
)

type EventBus struct {
	server element.Server
}

type Event = events.Event[map[string]any]

func New(server element.Server) *EventBus {
	return &EventBus{
		server: server,
	}
}

func (b *EventBus) Publish(ctx context.Context, dbKey string, event *Event) error {
	eventDao := dao.GetOutboxEventDao(dbKey)
	return events.Publish(ctx, eventDao, event)
}

func (b *EventBus) PublishList(ctx context.Context, dbKey string, eventList []*Event) error {
	eventDao := dao.GetOutboxEventDao(dbKey)
	return events.PublishList(ctx, eventDao, eventList)
}

func (b *EventBus) NewEvent(ctx context.Context, appId string, eventType string, data map[string]any, meta map[string]any) *Event {
	return NewEvent(ctx, appId, data, &events.EventOptions{Meta: meta, EventType: eventType})
}

func NewEvent(ctx context.Context, appId string, data map[string]any, opts ...*events.EventOptions) *Event {
	event := &Event{}
	event.SetData(ctx, appId, data, opts...)
	return event
}
