package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SuTaskDao struct {
	idao.Dao[*model.SuTask]
}

func NewSuTaskDao(dbKey string) *SuTaskDao {
	tableName := "analyse_task"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuTask{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuTask](newCfg)
	daoVal := &SuTaskDao{Dao: baseDao}
	return daoVal
}
