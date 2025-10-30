package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/case/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type CaseDao struct {
	idao.Dao[*model.Case]
}

func NewCaseDao(dbKey string) *CaseDao {
	tableName := "master_case"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Case{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Case](newCfg)
	daoVal := &CaseDao{Dao: baseDao}
	return daoVal
}
