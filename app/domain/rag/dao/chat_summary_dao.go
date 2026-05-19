package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type ChatSummaryDao struct {
	idao.Dao[*model.ChatSummary]
}

func NewChatSummaryDao(dbKey string) *ChatSummaryDao {
	tableName := "rag_chat_summary"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ChatSummary{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ChatSummary](newCfg)
	daoVal := &ChatSummaryDao{Dao: baseDao}
	return daoVal
}