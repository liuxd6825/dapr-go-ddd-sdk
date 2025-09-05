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

type FindByHomeIdQuery struct {
	HomeId string `json:"homeId" param:"homeId" required:"true"`
}

type FindByCodeQuery struct {
	Code   string `json:"code" param:"code" required:"true"`
	CaseId string `json:"caseId" param:"caseId" required:"-"`
	UserId string `json:"userId" param:"userId" required:"-"`
}

type FindCardFilesQuery struct {
	AppId string `json:"appId" param:"appId" required:"true"`
	FunId string `json:"funId" param:"funId" required:"-"`
}

type FindByAppIdQuery struct {
	AppId string `json:"appId" param:"appId" required:"true"`
}
