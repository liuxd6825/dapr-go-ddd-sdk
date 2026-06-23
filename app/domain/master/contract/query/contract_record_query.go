package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// ContractRecordFindByIdQuery 按ID查询
type ContractRecordFindByIdQuery = store.FindByIdRequest

// ContractRecordFindPagingQuery 分页查询
type ContractRecordFindPagingQuery = store.FindPagingQueryRequest

// ContractRecordFindByContractIdQuery 按合同ID查询
type ContractRecordFindByContractIdQuery struct {
	store.FindPagingQueryRequest
	ContractId string `json:"contractId" path:"contract-id" validate:"required" title:"合同ID"`
}