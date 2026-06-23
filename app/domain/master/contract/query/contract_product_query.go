package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ContractProductFindByIdQuery 按ID查询
type ContractProductFindByIdQuery = store.FindByIdRequest

// ContractProductFindPagingQuery 分页查询
type ContractProductFindPagingQuery = store.FindPagingQueryRequest

// ContractProductFindByContractIdQuery 按合同ID查询
type ContractProductFindByContractIdQuery struct {
	store.FindPagingQueryRequest
	ContractId string `json:"contractId" path:"contract-id" validate:"required" title:"合同ID"`
}