package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ContractHumanFindByIdQuery 按ID查询
type ContractHumanFindByIdQuery = store.FindByIdRequest

// ContractHumanFindPagingQuery 分页查询
type ContractHumanFindPagingQuery = store.FindPagingQueryRequest

// ContractHumanFindByContractIdQuery 按合同ID查询
type ContractHumanFindByContractIdQuery struct {
	store.FindPagingQueryRequest
	ContractId string `json:"contractId" path:"contract-id" validate:"required" title:"合同ID"`
}