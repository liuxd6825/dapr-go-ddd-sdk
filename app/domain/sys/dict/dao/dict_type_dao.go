package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type DictTypeDao struct {
	idao.Dao[*model.DictType]
}

func NewDictTypeDao(dbKey string) *DictTypeDao {
	tableName := "sys_dictionary_type"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.DictType{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.DictType](newCfg)
	daoVal := &DictTypeDao{Dao: baseDao}
	return daoVal
}
