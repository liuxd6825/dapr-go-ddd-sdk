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

// ProductRecordDao 产品-记录数据访问对象
// ⚠ 严格按源保留：model 中字段为 contractId（source bug），但级联删除使用 productId
type ProductRecordDao struct {
	idao.Dao[*model.ProductRecord]
}

func NewProductRecordDao(dbKey string) *ProductRecordDao {
	tableName := "product_record"
	dbSch := dbschema.NewDBSchemaWithStruct(tableName, &model.ProductRecord{}, tableName)
	newCfg := &dao.DaoConfig{
		DBKey:     dbKey,
		TableName: tableName,
		DBSchema:  dbSch,
	}
	baseDao := dao.NewDao[*model.ProductRecord](newCfg)
	return &ProductRecordDao{Dao: baseDao}
}

func (d *ProductRecordDao) FindPagingByProductId(ctx context.Context, qry store.FindPagingQuery, productId string, opts ...idao.CallOptions) store.FindPagingResult[*model.ProductRecord] {
	qry.SetMustFilter(fmt.Sprintf("productId=='%s'", productId))
	qry.SetPageSize(store.PagingMaxPageSize)
	return d.FindPaging(ctx, qry, opts...)
}

func (d *ProductRecordDao) DeleteByProductIds(ctx context.Context, productIds []string, opts ...idao.CallOptions) error {
	rsqlStr := rsql.NewBuilder().In("productId", productIds).Build()
	return d.DeleteByRSQL(ctx, rsqlStr, opts...).GetError()
}