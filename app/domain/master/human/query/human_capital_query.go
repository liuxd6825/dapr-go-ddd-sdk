package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanCapitalFindByIdQuery 按ID查询
type HumanCapitalFindByIdQuery = store.FindByIdRequest

// HumanCapitalFindPagingQuery 分页查询
type HumanCapitalFindPagingQuery = store.FindPagingQueryRequest

// HumanCapitalFindByHumanIdQuery 按人员ID查询
type HumanCapitalFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId"validate:"required" title:"人员ID"`
}
