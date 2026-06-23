package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ProductHumanFindByIdQuery 按ID查询
type ProductHumanFindByIdQuery = store.FindByIdRequest

// ProductHumanFindPagingQuery 分页查询
type ProductHumanFindPagingQuery = store.FindPagingQueryRequest

// ProductHumanFindByProductIdQuery 按产品ID查询
type ProductHumanFindByProductIdQuery struct {
	store.FindPagingQueryRequest
	ProductId string `json:"productId" path:"product-id" validate:"required" title:"产品ID"`
}