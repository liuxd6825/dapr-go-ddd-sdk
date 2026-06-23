package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
)

// ProductDao 产品数据访问对象
type ProductDao struct {
	idao.Dao[*model.Product]
}

func NewProductDao(dbKey string) *ProductDao {
	tableName := "product"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.Product{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.Product](newCfg)
	return &ProductDao{Dao: baseDao}
}

func (d *ProductDao) FindPagingByCaseId(ctx context.Context, qry store.FindPagingQuery, caseId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Product] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s'", caseId))
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductDao) FindPagingByTagId(ctx context.Context, qry store.FindPagingQuery, caseId, tagId string, opts ...idao.CallOptions) store.FindPagingResult[*model.Product] {
	qry.SetMustFilter(fmt.Sprintf("case_id=='%s' and tag_id=contains='%s'", caseId, tagId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductDao) CountByCaseId(ctx context.Context, caseId string, opts ...idao.CallOptions) (int64, error) {
	rsqlStr := fmt.Sprintf("case_id=='%s'", caseId)
	return d.CountByRSQL(ctx, rsqlStr, opts...)
}