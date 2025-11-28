package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type FolderMetaDao struct {
	idao.Dao[*model.FolderMeta]
}

func NewFolderMetaDao(dbKey string) *FolderMetaDao {
	tableName := "doc_folder_meta"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.FolderMeta{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.FolderMeta](newCfg)
	daoVal := &FolderMetaDao{Dao: baseDao}
	return daoVal
}
