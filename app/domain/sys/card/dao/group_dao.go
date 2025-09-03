package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type GroupDao struct {
	idao.Dao[*model.Group]
}

func NewGroupDao(dbKey string) *GroupDao {
	tableName := "sys_home_group"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Group{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Group](newCfg)
	daoVal := &GroupDao{Dao: baseDao}
	return daoVal
}
