package query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type FindByIdQuery struct {
	Id string `json:"id" param:"id" required:"true"`
}

type FindPagingByCaseIdQuery struct {
	CaseId string `json:"caseId"  param:"id" required:"true"`
	idao.FindPagingQueryRequest
}

type FindByCodeQuery struct {
	Code string `json:"code" path:"code" required:"true"`
}

type FindByDictTypeAndCodeQuery struct {
	Code     string `json:"code" path:"code" required:"true"`
	DictType string `json:"dictType" path:"dictType" required:"true"`
}

type FindByDictTypeCodeAndParamsQuery struct {
	DictTypeCode string `json:"dictTypeCode" path:"dictTypeCode" required:"true"`
	CaseTypeId   string `json:"caseTypeId" query:"caseTypeId" required:"-"`
	CaseId       string `json:"caseId" query:"caseId" required:"-"`
}
