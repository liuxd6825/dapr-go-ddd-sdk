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

// HumanSuspectAmountDao 嫌疑人金额数据访问对象
type HumanSuspectAmountDao struct {
	idao.Dao[*model.HumanSuspectAmount]
}

func NewHumanSuspectAmountDao(dbKey string) *HumanSuspectAmountDao {
	tableName := "human_suspect_amount"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.HumanSuspectAmount{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.HumanSuspectAmount](newCfg)
	return &HumanSuspectAmountDao{Dao: baseDao}
}

func (d *HumanSuspectAmountDao) FindPagingByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanSuspectAmount] {
	qry.SetMustFilter(fmt.Sprintf("humanId=='%s'", humanId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *HumanSuspectAmountDao) DeleteByHumanIds(ctx context.Context, humanIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("humanId", humanIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}