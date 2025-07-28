package events

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
)

func PublishEvent(ctx context.Context, dao idao.OutboxEventDao, appId string, data any, meta map[string]any) error {
	tenantId, err := appctx.GetTenantId3(ctx)
	if err != nil {
		return err
	}
	event, err := NewOutboxEvent(ctx, appId, tenantId, data, meta)
	if err != nil {
		return err
	}
	return dao.Create(ctx, event).GetError()
}

func PublishEvents(ctx context.Context, dao idao.OutboxEventDao, appId string, data []any, meta map[string]any) error {
	tenantId, err := appctx.GetTenantId3(ctx)
	if err != nil {
		return err
	}

	var events []*dbevent.OutboxEvent
	for _, data := range data {
		event, err := NewOutboxEvent(ctx, appId, tenantId, data, meta)
		if err != nil {
			return err
		}
		events = append(events, event)
	}
	return dao.CreateMany(ctx, events).GetError()
}
