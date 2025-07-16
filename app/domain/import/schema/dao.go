package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SchemaDao struct {
	idao.Dao[*SchemaModel]
}

func NewSchemaDao(dbKey string) *SchemaDao {
	tableName := "import_schema"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &SchemaModel{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*SchemaModel](newCfg)
	daoVal := &SchemaDao{Dao: baseDao}
	return daoVal
}
