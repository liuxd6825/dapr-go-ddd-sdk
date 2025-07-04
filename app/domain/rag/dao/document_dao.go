package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type DocumentDao struct {
	idao.Dao[*model.Document]
}

func NewDocumentDao(dbKey string) *DocumentDao {
	tableName := "rag_document"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Document{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Document](newCfg)
	daoVal := &DocumentDao{Dao: baseDao}
	return daoVal
}
