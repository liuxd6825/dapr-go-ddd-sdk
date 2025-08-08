package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type ExcelRowDao struct {
	idao.Dao[*model.ExcelRow]
}

func NewExcelRowDao(dbKey string) *ExcelRowDao {
	tableName := "import_excel_row"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ExcelRow{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ExcelRow](newCfg)
	daoVal := &ExcelRowDao{Dao: baseDao}
	return daoVal
}
