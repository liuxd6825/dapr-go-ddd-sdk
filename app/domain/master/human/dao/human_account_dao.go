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

// HumanAccountDao 人员-账户关联数据访问对象
type HumanAccountDao struct {
	idao.Dao[*model.HumanAccount]
}

func NewHumanAccountDao(dbKey string) *HumanAccountDao {
	tableName := "human_account"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.HumanAccount{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.HumanAccount](newCfg)
	return &HumanAccountDao{Dao: baseDao}
}

func (d *HumanAccountDao) FindPagingByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanAccount] {
	qry.SetMustFilter(fmt.Sprintf("humanId=='%s'", humanId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *HumanAccountDao) DeleteByHumanIds(ctx context.Context, humanIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("humanId", humanIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}