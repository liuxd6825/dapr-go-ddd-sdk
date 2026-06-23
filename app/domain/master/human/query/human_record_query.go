package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanRecordFindByIdQuery 按ID查询
type HumanRecordFindByIdQuery = store.FindByIdRequest

// HumanRecordFindPagingQuery 分页查询
type HumanRecordFindPagingQuery = store.FindPagingQueryRequest

// HumanRecordFindByHumanIdQuery 按人员ID查询
type HumanRecordFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
