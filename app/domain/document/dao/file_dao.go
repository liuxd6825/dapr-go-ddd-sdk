package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type FileDao struct {
	idao.Dao[*model.File]
}

func NewFileDao(dbKey string) *FileDao {
	tableName := "doc_file"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.File{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.File](newCfg)
	daoVal := &FileDao{Dao: baseDao}
	return daoVal
}
