package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ProductCompanyFindByIdQuery 按ID查询
type ProductCompanyFindByIdQuery = store.FindByIdRequest

// ProductCompanyFindPagingQuery 分页查询
type ProductCompanyFindPagingQuery = store.FindPagingQueryRequest

// ProductCompanyFindByProductIdQuery 按产品ID查询
type ProductCompanyFindByProductIdQuery struct {
	store.FindPagingQueryRequest
	ProductId string `json:"productId" path:"product-id" validate:"required" title:"产品ID"`
}