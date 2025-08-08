package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type ChatDao struct {
	idao.Dao[*model.Chat]
}

func NewChatDao(dbKey string) *ChatDao {
	tableName := "rag_chat"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Chat{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Chat](newCfg)
	daoVal := &ChatDao{Dao: baseDao}
	return daoVal
}
