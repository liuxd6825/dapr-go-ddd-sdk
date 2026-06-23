package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/contract/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

// ContractDao 合同数据访问对象
type ContractDao struct {
	idao.Dao[*model.Contract]
}

func NewContractDao(dbKey string) *ContractDao {
	tableName := "contract"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Contract{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Contract](newCfg)
	return &ContractDao{Dao: baseDao}
}

// FindPagingByCaseId 按 caseId 分页查询合同
func (d *ContractDao) FindPagingByCaseId(ctx context.Context, qry store.FindPagingQuery, caseId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Contract] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", caseId))
	return d.FindPaging(ctx, qry, opts...)
}

// FindPagingByTagId 按 caseId + tagId 分页查询合同
func (d *ContractDao) FindPagingByTagId(ctx context.Context, qry store.FindPagingQuery, caseId, tagId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Contract] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s' and tag_id=contains='%s'", caseId, tagId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// CountByCaseId 按 caseId 统计合同数量
func (d *ContractDao) CountByCaseId(ctx context.Context, caseId string, opts ...idao.CallOptions) (int64, error) {
	rsqlStr := fmt.Sprintf("case_id=='%s'", caseId)
	return d.CountByRSQL(ctx, rsqlStr, opts...)
}