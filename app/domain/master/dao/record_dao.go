package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type RecordDao struct {
	idao.Dao[*model.Record]
}

func NewRecordDao(dbKey string) *RecordDao {
	tableName := "master_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Record{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Record](newCfg)
	daoVal := &RecordDao{Dao: baseDao}
	return daoVal
}
