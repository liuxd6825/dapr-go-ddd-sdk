package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TagDao struct {
	idao.Dao[*model.Tag]
}

func NewTagDao(dbKey string) *TagDao {
	tableName := "sys_tag"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Tag{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.Tag](newCfg)
	daoVal := &TagDao{Dao: baseDao}
	return daoVal
}
