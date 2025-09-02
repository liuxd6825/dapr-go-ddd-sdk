package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type FunDao struct {
	idao.Dao[*model.Fun]
}

func NewFunDao(dbKey string) *FunDao {
	tableName := "sys_home_fun"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Fun{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Fun](newCfg)
	daoVal := &FunDao{Dao: baseDao}
	return daoVal
}
