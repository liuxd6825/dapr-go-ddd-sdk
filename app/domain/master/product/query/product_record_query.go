package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ProductRecordFindByIdQuery 按ID查询
type ProductRecordFindByIdQuery = store.FindByIdRequest

// ProductRecordFindPagingQuery 分页查询
type ProductRecordFindPagingQuery = store.FindPagingQueryRequest

// ProductRecordFindByProductIdQuery 按产品ID查询
type ProductRecordFindByProductIdQuery struct {
	store.FindPagingQueryRequest
	ProductId string `json:"productId" path:"product-id" validate:"required" title:"产品ID"`
}