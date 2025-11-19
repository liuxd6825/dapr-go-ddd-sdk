package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type CurrencyDao struct {
	idao.Dao[*model.Currency]
}

func NewCurrencyDao(dbKey string) *CurrencyDao {
	tableName := "sys_currency"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Currency{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Currency](newCfg)
	daoVal := &CurrencyDao{Dao: baseDao}
	return daoVal
}
