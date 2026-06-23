package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ProductContractFindByIdQuery 按ID查询
type ProductContractFindByIdQuery = store.FindByIdRequest

// ProductContractFindPagingQuery 分页查询
type ProductContractFindPagingQuery = store.FindPagingQueryRequest

// ProductContractFindByProductIdQuery 按产品ID查询
type ProductContractFindByProductIdQuery struct {
	store.FindPagingQueryRequest
	ProductId string `json:"productId" path:"product-id" validate:"required" title:"产品ID"`
}