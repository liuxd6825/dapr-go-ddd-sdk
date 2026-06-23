package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

type CompanyContractFindByIdQuery = store.FindByIdRequest
type CompanyContractFindPagingQuery = store.FindPagingQueryRequest

type CompanyContractFindByCompanyIdQuery struct {
	store.FindPagingQueryRequest
	CompanyId string `json:"companyId" param:"companyId"  validate:"required" title:"公司ID"`
}
