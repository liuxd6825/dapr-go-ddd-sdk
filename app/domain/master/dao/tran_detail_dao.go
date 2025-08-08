package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TranDetailDao struct {
	idao.Dao[*model.TranDetail]
}

func NewTranDetailDao(dbKey string) *TranDetailDao {
	tableName := "master_tran_detail"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.TranDetail{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.TranDetail](newCfg)
	daoVal := &TranDetailDao{Dao: baseDao}
	return daoVal
}
