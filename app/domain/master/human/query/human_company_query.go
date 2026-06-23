package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanCompanyFindByIdQuery 按ID查询
type HumanCompanyFindByIdQuery = store.FindByIdRequest

// HumanCompanyFindPagingQuery 分页查询
type HumanCompanyFindPagingQuery = store.FindPagingQueryRequest

// HumanCompanyFindByHumanIdQuery 按人员ID查询
type HumanCompanyFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
