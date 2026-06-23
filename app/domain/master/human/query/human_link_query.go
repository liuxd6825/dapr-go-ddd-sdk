package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanLinkFindByIdQuery 按ID查询
type HumanLinkFindByIdQuery = store.FindByIdRequest

// HumanLinkFindPagingQuery 分页查询
type HumanLinkFindPagingQuery = store.FindPagingQueryRequest

// HumanLinkFindByHumanIdQuery 按人员ID查询
type HumanLinkFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
