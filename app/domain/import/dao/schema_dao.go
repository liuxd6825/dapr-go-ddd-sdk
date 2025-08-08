package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SchemaDao struct {
	idao.Dao[*model.SchemaModel]
}

func NewSchemaDao(dbKey string) *SchemaDao {
	tableName := "import_schema"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SchemaModel{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SchemaModel](newCfg)
	daoVal := &SchemaDao{Dao: baseDao}
	return daoVal
}
