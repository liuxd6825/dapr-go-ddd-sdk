package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

type CompanyHumanFindByIdQuery = store.FindByIdRequest
type CompanyHumanFindPagingQuery = store.FindPagingQueryRequest

type CompanyHumanFindByCompanyIdQuery struct {
	store.FindPagingQueryRequest
	CompanyId string `json:"companyId" param:"companyId" validate:"required" title:"公司ID"`
}
