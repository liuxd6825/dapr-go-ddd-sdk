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

// HumanContractDao 人员-合同关联数据访问对象
type HumanContractDao struct {
	idao.Dao[*model.HumanContract]
}

func NewHumanContractDao(dbKey string) *HumanContractDao {
	tableName := "human_contract"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.HumanContract{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.HumanContract](newCfg)
	return &HumanContractDao{Dao: baseDao}
}

func (d *HumanContractDao) FindPagingByHumanId(ctx context.Context, qry store.FindPagingQuery, humanId string, opts ...idao.CallOptions) store.FindPagingResult[*model.HumanContract] {
	qry.SetMustFilter(fmt.Sprintf("humanId=='%s'", humanId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *HumanContractDao) DeleteByHumanIds(ctx context.Context, humanIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("humanId", humanIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}