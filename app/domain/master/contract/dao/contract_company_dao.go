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

// ContractCompanyDao 合同-公司关联数据访问对象
type ContractCompanyDao struct {
	idao.Dao[*model.ContractCompany]
}

func NewContractCompanyDao(dbKey string) *ContractCompanyDao {
	tableName := "contract_company"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ContractCompany{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ContractCompany](newCfg)
	return &ContractCompanyDao{Dao: baseDao}
}

// FindPagingByContractId 按合同ID分页查询
func (d *ContractCompanyDao) FindPagingByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractCompany] {
	qry.SetMustFilter(fmt.Sprintf("contractId=='%s'", contractId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// DeleteByContractIds 级联按合同IDs删除
func (d *ContractCompanyDao) DeleteByContractIds(ctx context.Context, contractIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("contractId", contractIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}