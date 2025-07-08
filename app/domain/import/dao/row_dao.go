package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type RowDao struct {
	idao.Dao[*model.Row]
}

func NewRowDao(dbKey string) *RowDao {
	tableName := "import_row"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Row{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Row](newCfg)
	daoVal := &RowDao{Dao: baseDao}
	return daoVal
}
