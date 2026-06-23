package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanAddressFindByIdQuery 按ID查询
type HumanAddressFindByIdQuery = store.FindByIdRequest

// HumanAddressFindPagingQuery 分页查询
type HumanAddressFindPagingQuery = store.FindPagingQueryRequest

// HumanAddressFindByHumanIdQuery 按人员ID查询
type HumanAddressFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
