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

// ProductContractDao 产品-合同关联数据访问对象
type ProductContractDao struct {
	idao.Dao[*model.ProductContract]
}

func NewProductContractDao(dbKey string) *ProductContractDao {
	tableName := "product_contract"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ProductContract{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ProductContract](newCfg)
	return &ProductContractDao{Dao: baseDao}
}

func (d *ProductContractDao) FindPagingByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductContract] {
	qry.SetMustFilter(fmt.Sprintf("productId=='%s'", productId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductContractDao) DeleteByProductIds(ctx context.Context, productIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("productId", productIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}