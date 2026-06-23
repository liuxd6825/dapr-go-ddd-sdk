package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ContractContractFindByIdQuery 按ID查询
type ContractContractFindByIdQuery = store.FindByIdRequest

// ContractContractFindPagingQuery 分页查询
type ContractContractFindPagingQuery = store.FindPagingQueryRequest

// ContractContractFindByContractIdQuery 按合同ID查询
type ContractContractFindByContractIdQuery struct {
	store.FindPagingQueryRequest
	ContractId string `json:"contractId" path:"contract-id" validate:"required" title:"合同ID"`
}