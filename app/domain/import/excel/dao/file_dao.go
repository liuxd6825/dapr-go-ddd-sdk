package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type FileDao struct {
	idao.Dao[*model.ExcelFile]
}

func NewFileDao(dbKey string) *FileDao {
	tableName := "import_excel_file"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ExcelFile{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.ExcelFile](newCfg)
	daoVal := &FileDao{Dao: baseDao}
	return daoVal
}
