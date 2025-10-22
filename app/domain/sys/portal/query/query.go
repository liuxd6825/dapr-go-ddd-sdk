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
