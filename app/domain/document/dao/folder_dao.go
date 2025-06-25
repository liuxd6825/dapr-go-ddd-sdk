package dao

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type FolderDao struct {
	idao.Dao[*model.Folder]
}

func NewFolderDao(dbKey string) *FolderDao {
	tableName := "doc_folder"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Folder{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Folder](newCfg)
	daoVal := &FolderDao{Dao: baseDao}
	daoVal.Table().AutoMigrate(context.Background())
	return daoVal
}
