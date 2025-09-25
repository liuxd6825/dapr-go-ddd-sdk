package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SuTaskLogDao struct {
	idao.Dao[*model.TaskLog]
}

func NewSuTaskLogDao(dbKey string) *SuTaskLogDao {
	tableName := "su_task_log"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.TaskLog{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.TaskLog](newCfg)
	daoVal := &SuTaskLogDao{Dao: baseDao}
	return daoVal
}
