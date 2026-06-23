package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanAccountFindByIdQuery 按ID查询
type HumanAccountFindByIdQuery = store.FindByIdRequest

// HumanAccountFindPagingQuery 分页查询
type HumanAccountFindPagingQuery = store.FindPagingQueryRequest

// HumanAccountFindByHumanIdQuery 按人员ID查询
type HumanAccountFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
