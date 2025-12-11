package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/status/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type StatusDao struct {
	idao.Dao[*model.Status]
}

func NewStatusDao(dbKey string) *StatusDao {
	tableName := "sys_status"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Status{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Status](newCfg)
	daoVal := &StatusDao{Dao: baseDao}
	return daoVal
}
