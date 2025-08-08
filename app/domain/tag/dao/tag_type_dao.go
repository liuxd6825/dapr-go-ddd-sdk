package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

type TagTypeDao struct {
	idao.Dao[*model.TagType]
}

func NewTagTypeDao(dbKey string) *TagTypeDao {
	tableName := "sys_tag_type"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.TagType{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.TagType](newCfg)
	daoVal := &TagTypeDao{Dao: baseDao}
	return daoVal
}
