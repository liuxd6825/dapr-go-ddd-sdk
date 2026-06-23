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

// ContractRecordDao 合同-记录关联数据访问对象
type ContractRecordDao struct {
	idao.Dao[*model.ContractRecord]
}

func NewContractRecordDao(dbKey string) *ContractRecordDao {
	tableName := "contract_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ContractRecord{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ContractRecord](newCfg)
	return &ContractRecordDao{Dao: baseDao}
}

// FindPagingByContractId 按合同ID分页查询
func (d *ContractRecordDao) FindPagingByContractId(ctx context.Context, qry store.FindPagingQuery, contractId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ContractRecord] {
	qry.SetMustFilter(fmt.Sprintf("contractId=='%s'", contractId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}