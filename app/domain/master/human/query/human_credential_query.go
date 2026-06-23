package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// HumanCredentialFindByIdQuery 按ID查询
type HumanCredentialFindByIdQuery = store.FindByIdRequest

// HumanCredentialFindPagingQuery 分页查询
type HumanCredentialFindPagingQuery = store.FindPagingQueryRequest

// HumanCredentialFindByHumanIdQuery 按人员ID查询
type HumanCredentialFindByHumanIdQuery struct {
	store.FindPagingQueryRequest
	HumanId string `json:"humanId" path:"humanId" validate:"required" title:"人员ID"`
}
