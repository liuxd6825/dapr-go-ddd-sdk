package dao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/product/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/rsql"
)

// ProductCompanyDao 产品-公司关联数据访问对象
type ProductCompanyDao struct {
	idao.Dao[*model.ProductCompany]
}

func NewProductCompanyDao(dbKey string) *ProductCompanyDao {
	tableName := "product_company"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ProductCompany{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ProductCompany](newCfg)
	return &ProductCompanyDao{Dao: baseDao}
}

func (d *ProductCompanyDao) FindPagingByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductCompany] {
	qry.SetMustFilter(fmt.Sprintf("productId=='%s'", productId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductCompanyDao) DeleteByProductIds(ctx context.Context, productIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("productId", productIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}