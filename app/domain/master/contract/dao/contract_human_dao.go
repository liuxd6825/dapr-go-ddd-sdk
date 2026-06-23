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

// ContractHumanDao 合同-人员关联数据访问对象
type ContractHumanDao struct {
	idao.Dao[*model.ContractHuman]
}

func NewContractHumanDao(dbKey string) *ContractHumanDao {
	tableName := "contract_human"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ContractHuman{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ContractHuman](newCfg)
	return &ContractHumanDao{Dao: baseDao}
}

// FindPagingByContractId 按合同ID分页查询
func (d *ContractHumanDao) FindPagingByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractHuman] {
	qry.SetMustFilter(fmt.Sprintf("contractId=='%s'", contractId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// DeleteByContractIds 级联按合同IDs删除
func (d *ContractHumanDao) DeleteByContractIds(ctx context.Context, contractIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("contractId", contractIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}