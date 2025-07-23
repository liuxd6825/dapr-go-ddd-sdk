package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type OutboxEventDao struct {
	idao.Dao[*dbevent.OutboxEvent]
}

func NewOutboxEventDao(dbKey string) idao.OutboxEventDao {
	tableName := "sys_outbox_event"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &dbevent.OutboxEvent{}, tableName)
	newCfg := &NewConfig{
		DBKey:      dbKey,
		IsPubEvent: IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := NewDao[*dbevent.OutboxEvent](newCfg)
	dao := &OutboxEventDao{Dao: baseDao}
	//dao.Table().AutoMigrate(context.Background())
	return dao
}
