package query

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"

type CompanyRecordFindByIdQuery = store.FindByIdRequest
type CompanyRecordFindPagingQuery = store.FindPagingQueryRequest

type CompanyRecordFindByCompanyIdQuery struct {
	store.FindPagingQueryRequest
	CompanyId string `json:"companyId" query:"companyId" validate:"required" title:"公司ID"`
}
