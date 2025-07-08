package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type FileDao struct {
	idao.Dao[*model.File]
}

func NewFileDao(dbKey string) *FileDao {
	tableName := "excel_file"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.File{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.File](newCfg)
	daoVal := &FileDao{Dao: baseDao}
	return daoVal
}
