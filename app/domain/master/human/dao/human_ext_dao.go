package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

// HumanExtDao 人员-扩展信息数据访问对象
type HumanExtDao struct {
	idao.Dao[*model.HumanExt]
}

func NewHumanExtDao(dbKey string) *HumanExtDao {
	tableName := "human_ext"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.HumanExt{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.HumanExt](newCfg)
	return &HumanExtDao{Dao: baseDao}
}

func (d *HumanExtDao) FindPagingByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanExt] {
	qry.SetMustFilter(fmt.Sprintf("humanId=='%s'", humanId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *HumanExtDao) DeleteByHumanIds(ctx context.Context, humanIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("humanId", humanIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}