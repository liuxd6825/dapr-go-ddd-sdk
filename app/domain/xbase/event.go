package xbase

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"sync"
)

type Event[T any] struct {
	EventId    string      `json:"eventId" bson:"eventId" validate:"required" `
	OccurredOn *times.Time `json:"occurredOn" bson:"occurredOn" validate:"required" `
	Data       T           `json:"data"  bson:"data" validate:"required"`
}

var _outboxDao idao.OutboxEventDao
var _outboxOnce sync.Once

func newOutboxEventDao() idao.OutboxEventDao {
	_outboxOnce.Do(func() {
		_outboxDao = dao.NewOutboxEventDao(config.DBKey)
	})
	return _outboxDao
}

func PublishEvent(ctx context.Context, appId string, data any, meta map[string]any) error {
	outboxEventDao := newOutboxEventDao()
	return events.PublishEvent(ctx, outboxEventDao, appId, data, meta)
}

func PublishEvents(ctx context.Context, appId string, data []any, meta map[string]any) error {
	outboxEventDao := newOutboxEventDao()
	return events.PublishEvents(ctx, outboxEventDao, appId, data, meta)
}
