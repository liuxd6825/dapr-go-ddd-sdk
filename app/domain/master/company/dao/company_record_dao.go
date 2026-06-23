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

// CompanyRecordDao 公司-流水关联数据访问对象（只读）
type CompanyRecordDao struct {
	idao.Dao[*model.CompanyRecord]
}

func NewCompanyRecordDao(dbKey string) *CompanyRecordDao {
	tableName := "company_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.CompanyRecord{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.CompanyRecord](newCfg)
	return &CompanyRecordDao{Dao: baseDao}
}

// FindPagingByCompanyId 按公司ID分页查询
func (d *CompanyRecordDao) FindPagingByCompanyId(ctx context.Context, qry store.FindPagingQuery, companyId string, opts ...idao.CallOptions) store.FindPagingResult[*model.CompanyRecord] {
	qry.SetMustFilter(fmt.Sprintf("companyId=='%s'", companyId))
	return d.FindPaging(ctx, qry, opts...)
}

// DeleteByCompanyIds 级联按公司IDs删除
func (d *CompanyRecordDao) DeleteByCompanyIds(ctx context.Context, companyIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("companyId", companyIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}