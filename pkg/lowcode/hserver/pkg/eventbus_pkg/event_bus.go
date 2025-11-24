package eventbus_pkg

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
)

type EventBus struct {
	server element.Server
}

func New(server element.Server) *EventBus {
	return &EventBus{
		server: server,
	}
}

func (b *EventBus) PublishEvent(ctx context.Context, dbKey string, appId string, data any, meta map[string]any) error {
	eventDao := dao.GetOutboxEventDao(dbKey)
	return events.PublishEvent(ctx, eventDao, appId, data, meta)
}

func (b *EventBus) PublishEvents(ctx context.Context, dbKey string, appId string, data []any, meta map[string]any) error {
	eventDao := dao.GetOutboxEventDao(dbKey)
	return events.PublishEvents(ctx, eventDao, appId, data, meta)
}
