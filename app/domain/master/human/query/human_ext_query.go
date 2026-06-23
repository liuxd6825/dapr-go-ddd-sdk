package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanExtFindByIdQuery 按ID查询
type HumanExtFindByIdQuery = store.FindByIdRequest

// HumanExtFindPagingQuery 分页查询
type HumanExtFindPagingQuery = store.FindPagingQueryRequest

// HumanExtFindByHumanIdQuery 按人员ID查询
type HumanExtFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
