package dao

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
)

type OutboxEventDao struct {
	idao.Dao[*events.EventEntity]
}

var daoCaches = types.NewCMap[events.EventPublisher]()

func NewOutboxEventDao(dbKey string) events.EventPublisher {
	tableName := "sys_outbox_event"
	baseDao := NewDao[*events.EventEntity](NewConfig(dbKey, tableName, &events.EventEntity{}))
	dao := &OutboxEventDao{Dao: baseDao}
	return dao
}

func GetOutboxEventDao(dbKey string) events.EventPublisher {
	dao, ok := daoCaches.Get(dbKey)
	if ok {
		return dao
	}
	dao = NewOutboxEventDao(dbKey)
	daoCaches.Set(dbKey, dao)
	return dao
}

func (d *OutboxEventDao) Publish(ctx context.Context, event *events.EventEntity) error {
	return d.Dao.Create(ctx, event).GetError()
}

func (d *OutboxEventDao) PublishList(ctx context.Context, data []*events.EventEntity) error {
	return d.Dao.CreateMany(ctx, data).GetError()
}
