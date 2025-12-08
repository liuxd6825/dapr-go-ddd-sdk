package query

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
)

type FindByIdQuery struct {
	Id string `json:"id" param:"id" required:"true"`
}

type FindPagingByTenantIdRequest struct {
	TenantId string `json:"tenantId"  query:"tenant-id" required:"true"`
	idao.FindPagingQueryRequest
}

type FindByUserIdQuery struct {
	UserId string `json:"userId" query:"user-id" required:"true"`
}

type FindByTenantIdQuery struct {
	TenantId string `json:"tenantId" query:"tenant-id" required:"true"`
}
