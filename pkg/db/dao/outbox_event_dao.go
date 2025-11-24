package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
)

type OutboxEventDao struct {
	idao.Dao[*dbevent.OutboxEvent]
}

var daoCaches = types.NewCMap[idao.OutboxEventDao]()

func NewOutboxEventDao(dbKey string) idao.OutboxEventDao {
	tableName := "sys_outbox_event"
	baseDao := NewDao[*dbevent.OutboxEvent](NewConfig(dbKey, tableName, &dbevent.OutboxEvent{}))
	dao := &OutboxEventDao{Dao: baseDao}
	return dao
}

func GetOutboxEventDao(dbKey string) idao.OutboxEventDao {
	dao, ok := daoCaches.Get(dbKey)
	if ok {
		return dao
	}
	dao = NewOutboxEventDao(dbKey)
	daoCaches.Set(dbKey, dao)
	return dao
}
