package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type TagDao struct {
	idao.Dao[*model.Tag]
}

func NewTagDao(dbKey string) *TagDao {
	tableName := "sys_tag"
	baseDao := dao.NewDao[*model.Tag](dao.NewConfig(dbKey, tableName, &model.Tag{}))
	daoVal := &TagDao{Dao: baseDao}
	return daoVal
}
