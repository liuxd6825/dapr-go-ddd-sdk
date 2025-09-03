package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type CardFileDao struct {
	idao.Dao[*model.CardFile]
}

func NewCardFileDao(dbKey string) *CardFileDao {
	tableName := "sys_home_card_file"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CardFile{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CardFile](newCfg)
	daoVal := &CardFileDao{Dao: baseDao}
	return daoVal
}
