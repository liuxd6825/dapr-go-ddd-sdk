package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

type CompanyCompanyFindByIdQuery = store.FindByIdRequest
type CompanyCompanyFindPagingQuery = store.FindPagingQueryRequest

type CompanyCompanyFindByCompanyIdQuery struct {
	store.FindPagingQueryRequest
	CompanyId string `json:"companyId"  param:"companyId"  validate:"required" title:"公司ID"`
}
