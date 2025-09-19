package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SuBatchDao struct {
	idao.Dao[*model.SuBatch]
}

func NewSuBatchDao(dbKey string) *SuBatchDao {
	tableName := "su_batch"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuBatch{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuBatch](newCfg)
	daoVal := &SuBatchDao{Dao: baseDao}
	return daoVal
}
