package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type RecordDao struct {
	idao.Dao[*model.Record]
}

func NewRecordDao(dbKey string) *RecordDao {
	tableName := "import_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Record{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Record](newCfg)
	daoVal := &RecordDao{Dao: baseDao}
	return daoVal
}
