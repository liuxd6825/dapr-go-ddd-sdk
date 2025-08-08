package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
)

type OutboxEventDao struct {
	idao.Dao[*dbevent.OutboxEvent]
}

func NewOutboxEventDao(dbKey string) idao.OutboxEventDao {
	tableName := "sys_outbox_event"
	baseDao := NewDao[*dbevent.OutboxEvent](NewConfig(dbKey, tableName, &dbevent.OutboxEvent{}))
	dao := &OutboxEventDao{Dao: baseDao}
	return dao
}
