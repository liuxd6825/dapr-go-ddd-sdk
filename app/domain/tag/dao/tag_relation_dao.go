package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TagRelationDao struct {
	idao.Dao[*model.TagRelation]
}

func NewTagRelationDao(dbKey string) *TagRelationDao {
	tableName := "sys_tag_relation"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.TagRelation{}, tableName)
	newCfg := &dao.NewConfig{
		DBKey:      dbKey,
		IsPubEvent: dao.IsFalse(),
		TableName:  tableName,
		DBSchema:   dbSch,
	}
	baseDao := dao.NewDao[*model.TagRelation](newCfg)
	daoVal := &TagRelationDao{Dao: baseDao}
	return daoVal
}
