package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SuBatchItemDao struct {
	idao.Dao[*model.SuBatchItem]
}

func NewSuBatchItemDao(dbKey string) *SuBatchItemDao {
	tableName := "su_batch_item"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuBatchItem{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuBatchItem](newCfg)
	daoVal := &SuBatchItemDao{Dao: baseDao}
	return daoVal
}
