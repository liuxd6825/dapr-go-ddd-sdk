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

// CompanyDao 公司数据访问对象
type CompanyDao struct {
	idao.Dao[*model.Company]
}

func NewCompanyDao(dbKey string) *CompanyDao {
	tableName := "company"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Company{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Company](newCfg)
	return &CompanyDao{Dao: baseDao}
}

// FindPagingByCaseId 按 caseId 分页查询公司
func (d *CompanyDao) FindPagingByCaseId(ctx context.Context, qry store.FindPagingQuery, caseId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Company] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", caseId))
	return d.FindPaging(ctx, qry, opts...)
}

// FindPagingByTagId 按 caseId + tagId 分页查询公司
func (d *CompanyDao) FindPagingByTagId(ctx context.Context, qry store.FindPagingQuery, caseId, tagId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Company] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s' and tag_id=contains='%s'", caseId, tagId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

// CountByCaseId 按 caseId 统计公司数量
func (d *CompanyDao) CountByCaseId(ctx context.Context, caseId string, opts ...idao.CallOptions) (int64, error) {
	rsqlStr := rsql.NewBuilder().Eq("case_id", caseId).Build()
	return d.CountByRSQL(ctx, rsqlStr, opts...)
}