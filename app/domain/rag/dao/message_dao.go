package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type MessageDao[T interface{ *model.Message }] struct {
	idao.Dao[*model.Message]
}

func NewMessageDao(dbKey string) idao.Dao[*model.Message] {
	tableName := "rag_message"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Message{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Message](newCfg)
	daoVal := &MessageDao[*model.Message]{Dao: baseDao}
	daoVal.Table().AutoMigrate(context.Background())
	return daoVal
}
