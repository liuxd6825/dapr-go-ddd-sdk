package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanContractFindByIdQuery 按ID查询
type HumanContractFindByIdQuery = store.FindByIdRequest

// HumanContractFindPagingQuery 分页查询
type HumanContractFindPagingQuery = store.FindPagingQueryRequest

// HumanContractFindByHumanIdQuery 按人员ID查询
type HumanContractFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
