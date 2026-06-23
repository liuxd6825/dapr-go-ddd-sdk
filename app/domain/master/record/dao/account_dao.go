package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type AccountDao struct {
	idao.Dao[*model.Account]
}

func NewAccountDao(dbKey string) *AccountDao {
	tableName := "master_account"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Account{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Account](newCfg)
	daoVal := &AccountDao{Dao: baseDao}
	return daoVal
}
