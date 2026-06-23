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

// ProductProductDao 产品-产品关联数据访问对象
type ProductProductDao struct {
	idao.Dao[*model.ProductProduct]
}

func NewProductProductDao(dbKey string) *ProductProductDao {
	tableName := "product_product"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ProductProduct{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ProductProduct](newCfg)
	return &ProductProductDao{Dao: baseDao}
}

func (d *ProductProductDao) FindPagingByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductProduct] {
	qry.SetMustFilter(fmt.Sprintf("productId=='%s'", productId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductProductDao) DeleteByProductIds(ctx context.Context, productIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("productId", productIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}