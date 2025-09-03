package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type HomeDao struct {
	idao.Dao[*model.Home]
}

func NewHomeDao(dbKey string) *HomeDao {
	tableName := "sys_home_home"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Home{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Home](newCfg)
	daoVal := &HomeDao{Dao: baseDao}
	return daoVal
}
