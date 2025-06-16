package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type MessageDao struct {
	idao.Dao[*model.Message]
}

func NewMessageDao(dbKey string) *MessageDao {
	tableName := "rag_message"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Message{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Message](newCfg)
	daoVal := &MessageDao{Dao: baseDao}
	daoVal.Table().AutoMigrate(context.Background())
	return daoVal
}
