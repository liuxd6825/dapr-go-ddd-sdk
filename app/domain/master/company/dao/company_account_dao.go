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

// CompanyAccountDao 公司账号关联数据访问对象
type CompanyAccountDao struct {
	idao.Dao[*model.CompanyAccount]
}

func NewCompanyAccountDao(dbKey string) *CompanyAccountDao {
	tableName := "company_account"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CompanyAccount{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CompanyAccount](newCfg)
	return &CompanyAccountDao{Dao: baseDao}
}

// FindPagingByCompanyId 按公司ID分页查询
func (d *CompanyAccountDao) FindPagingByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyAccount] {
	qry.SetMustFilter(fmt.Sprintf("companyId=='%s'", companyId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// DeleteByCompanyIds 级联按公司IDs删除
func (d *CompanyAccountDao) DeleteByCompanyIds(ctx context.Context, companyIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("companyId", companyIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}