package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ProductProductFindByIdQuery 按ID查询
type ProductProductFindByIdQuery = store.FindByIdRequest

// ProductProductFindPagingQuery 分页查询
type ProductProductFindPagingQuery = store.FindPagingQueryRequest

// ProductProductFindByProductIdQuery 按产品ID查询
type ProductProductFindByProductIdQuery struct {
	store.FindPagingQueryRequest
	ProductId string `json:"productId" path:"product-id" validate:"required" title:"产品ID"`
}