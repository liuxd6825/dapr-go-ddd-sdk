package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type CardDao struct {
	idao.Dao[*model.Card]
}

func NewCardDao(dbKey string) *CardDao {
	tableName := "sys_home_card"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Card{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Card](newCfg)
	daoVal := &CardDao{Dao: baseDao}
	return daoVal
}
