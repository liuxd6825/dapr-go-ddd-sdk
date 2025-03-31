package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type OutboxDao[T interface{ *dbevent.Outbox }] struct {
	idao.Dao[*dbevent.Outbox]
}

func NewOutboxDao(dbKey string) idao.OutboxDao {
	tableName := "sys_outbox"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &dbevent.Outbox{}, tableName)
	newCfg := &NewConfig{
		DBKey:      dbKey,
		IsPubEvent: IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := NewDao[*dbevent.Outbox](newCfg)
	dao := &OutboxDao[*dbevent.Outbox]{Dao: baseDao}
	dao.Table().AutoMigrate(context.Background())
	return dao
}
