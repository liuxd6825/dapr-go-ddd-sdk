package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ContractCompanyFindByIdQuery 按ID查询
type ContractCompanyFindByIdQuery = store.FindByIdRequest

// ContractCompanyFindPagingQuery 分页查询
type ContractCompanyFindPagingQuery = store.FindPagingQueryRequest

// ContractCompanyFindByContractIdQuery 按合同ID查询
type ContractCompanyFindByContractIdQuery struct {
	store.FindPagingQueryRequest
	ContractId string `json:"contractId" path:"contract-id" validate:"required" title:"合同ID"`
}