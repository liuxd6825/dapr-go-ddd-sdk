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

// HumanHumanDao 人员-人员关联数据访问对象
type HumanHumanDao struct {
	idao.Dao[*model.HumanHuman]
}

func NewHumanHumanDao(dbKey string) *HumanHumanDao {
	tableName := "human_human"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.HumanHuman{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.HumanHuman](newCfg)
	return &HumanHumanDao{Dao: baseDao}
}

func (d *HumanHumanDao) FindPagingByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanHuman] {
	qry.SetMustFilter(fmt.Sprintf("humanId=='%s'", humanId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *HumanHumanDao) DeleteByHumanIds(ctx context.Context, humanIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("humanId", humanIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}