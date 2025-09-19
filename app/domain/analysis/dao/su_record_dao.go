package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type SuRecordDao struct {
	idao.Dao[*model.SuRecord]
}

func NewSuRecordDao(dbKey string) *SuRecordDao {
	tableName := "su_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.SuRecord{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.SuRecord](newCfg)
	daoVal := &SuRecordDao{Dao: baseDao}
	return daoVal
}
