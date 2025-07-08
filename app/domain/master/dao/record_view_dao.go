package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/view"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type RecordDao struct {
	idao.Dao[*view.RecordView]
}

func NewRecordDao(dbKey string) *RecordDao {
	tableName := "master_record_view"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &view.RecordView{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*view.RecordView](newCfg)
	daoVal := &RecordDao{Dao: baseDao}
	return daoVal
}
