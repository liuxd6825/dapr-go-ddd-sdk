package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type ChatDao[T interface{ *model.Chat }] struct {
	idao.Dao[*model.Chat]
}

func NewChatDao(dbKey string) idao.Dao[*model.Chat] {
	tableName := "rag_chat"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Chat{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Chat](newCfg)
	daoVal := &ChatDao[*model.Chat]{Dao: baseDao}
	daoVal.Table().AutoMigrate(context.Background())
	return daoVal
}
