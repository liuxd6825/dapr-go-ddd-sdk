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

// ProductHumanDao 产品-人员关联数据访问对象
type ProductHumanDao struct {
	idao.Dao[*model.ProductHuman]
}

func NewProductHumanDao(dbKey string) *ProductHumanDao {
	tableName := "product_human"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ProductHuman{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ProductHuman](newCfg)
	return &ProductHumanDao{Dao: baseDao}
}

func (d *ProductHumanDao) FindPagingByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductHuman] {
	qry.SetMustFilter(fmt.Sprintf("productId=='%s'", productId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductHumanDao) DeleteByProductIds(ctx context.Context, productIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("productId", productIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}