package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

// CompanyAccountFindByIdQuery 按ID查询
type CompanyAccountFindByIdQuery = store.FindByIdRequest

// CompanyAccountFindPagingQuery 分页查询
type CompanyAccountFindPagingQuery = store.FindPagingQueryRequest

// CompanyAccountFindByCompanyIdQuery 按公司ID查询
type CompanyAccountFindByCompanyIdQuery struct {
	store.FindPagingQueryRequest
	CompanyId string `json:"companyId" param:"companyId" validate:"required" title:"公司ID"`
}
