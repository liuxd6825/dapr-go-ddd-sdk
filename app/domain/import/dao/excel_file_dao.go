package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type ExcelFileDao struct {
	idao.Dao[*model.ExcelFile]
}

func NewExcelFileDao(dbKey string) *ExcelFileDao {
	tableName := "import_excel_file"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ExcelFile{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ExcelFile](newCfg)
	daoVal := &ExcelFileDao{Dao: baseDao}
	return daoVal
}
