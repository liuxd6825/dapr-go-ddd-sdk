package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/dict/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type DictDao struct {
	idao.Dao[*model.Dict]
}

func NewDictDao(dbKey string) *DictDao {
	tableName := "sys_dictionary"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Dict{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Dict](newCfg)
	daoVal := &DictDao{Dao: baseDao}
	return daoVal
}
