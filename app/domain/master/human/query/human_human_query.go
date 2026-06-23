package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanHumanFindByIdQuery 按ID查询
type HumanHumanFindByIdQuery = store.FindByIdRequest

// HumanHumanFindPagingQuery 分页查询
type HumanHumanFindPagingQuery = store.FindPagingQueryRequest

// HumanHumanFindByHumanIdQuery 按人员ID查询
type HumanHumanFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
