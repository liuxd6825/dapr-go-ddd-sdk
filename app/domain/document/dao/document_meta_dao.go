package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type DocumentMetaDao struct {
	idao.Dao[*model.DocumentMeta]
}

func NewDocumentMetaDao(dbKey string) *DocumentMetaDao {
	tableName := "doc_document_meta"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.DocumentMeta{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.DocumentMeta](newCfg)
	daoVal := &DocumentMetaDao{Dao: baseDao}
	return daoVal
}
