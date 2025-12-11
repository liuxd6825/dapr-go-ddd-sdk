package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type BankDao struct {
	idao.Dao[*model.Bank]
}

func NewBankDao(dbKey string) *BankDao {
	tableName := "sys_bank"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Bank{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Bank](newCfg)
	daoVal := &BankDao{Dao: baseDao}
	return daoVal
}
