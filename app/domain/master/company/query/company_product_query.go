package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

type CompanyProductFindByIdQuery = store.FindByIdRequest
type CompanyProductFindPagingQuery = store.FindPagingQueryRequest

type CompanyProductFindByCompanyIdQuery struct {
	store.FindPagingQueryRequest
	CompanyId string `json:"companyId"  param:"companyId"  validate:"required" title:"公司ID"`
}
