package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanProductFindByIdQuery 按ID查询
type HumanProductFindByIdQuery = store.FindByIdRequest

// HumanProductFindPagingQuery 分页查询
type HumanProductFindPagingQuery = store.FindPagingQueryRequest

// HumanProductFindByHumanIdQuery 按人员ID查询
type HumanProductFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
