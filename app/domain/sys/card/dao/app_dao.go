package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type AppDao struct {
	idao.Dao[*model.App]
}

func NewAppDao(dbKey string) *AppDao {
	tableName := "sys_home_app"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.App{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.App](newCfg)
	daoVal := &AppDao{Dao: baseDao}
	return daoVal
}
