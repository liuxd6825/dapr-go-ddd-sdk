package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

// ContractProductDao 合同-产品关联数据访问对象
type ContractProductDao struct {
	idao.Dao[*model.ContractProduct]
}

func NewContractProductDao(dbKey string) *ContractProductDao {
	tableName := "contract_product"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ContractProduct{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ContractProduct](newCfg)
	return &ContractProductDao{Dao: baseDao}
}

// FindPagingByContractId 按合同ID分页查询
func (d *ContractProductDao) FindPagingByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractProduct] {
	qry.SetMustFilter(fmt.Sprintf("contractId=='%s'", contractId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// DeleteByContractIds 级联按合同IDs删除
func (d *ContractProductDao) DeleteByContractIds(ctx context.Context, contractIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("contractId", contractIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}