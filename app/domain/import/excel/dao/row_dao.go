package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/excel/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type RowDao struct {
	idao.Dao[*model.ExcelRow]
}

func NewRowDao(dbKey string) *RowDao {
	tableName := "import_excel_row"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ExcelRow{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.ExcelRow](newCfg)
	daoVal := &RowDao{Dao: baseDao}
	return daoVal
}
