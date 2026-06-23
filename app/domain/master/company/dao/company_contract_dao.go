package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

type CompanyContractDao struct {
	idao.Dao[*model.CompanyContract]
}

func NewCompanyContractDao(dbKey string) *CompanyContractDao {
	tableName := "company_contract"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CompanyContract{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CompanyContract](newCfg)
	return &CompanyContractDao{Dao: baseDao}
}

func (d *CompanyContractDao) FindPagingByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyContract] {
	qry.SetMustFilter(fmt.Sprintf("companyId=='%s'", companyId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *CompanyContractDao) DeleteByCompanyIds(ctx context.Context, companyIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("companyId", companyIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}