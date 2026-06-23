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

type CompanyProductDao struct {
	idao.Dao[*model.CompanyProduct]
}

func NewCompanyProductDao(dbKey string) *CompanyProductDao {
	tableName := "company_product"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CompanyProduct{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CompanyProduct](newCfg)
	return &CompanyProductDao{Dao: baseDao}
}

func (d *CompanyProductDao) FindPagingByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyProduct] {
	qry.SetMustFilter(fmt.Sprintf("companyId=='%s'", companyId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *CompanyProductDao) DeleteByCompanyIds(ctx context.Context, companyIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("companyId", companyIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}