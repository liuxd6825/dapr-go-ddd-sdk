package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type TenantCreateCommand struct {
	xbase.Command[model.Tenant]
}
type TenantUpdateCommand struct {
	xbase.Command[model.Tenant]
}
type TenantDeleteCommand struct {
	xbase.DeleteByIdCommand
}

type InitTenant struct {
	TenantName string `json:"tenantName"`
	TenantId   string `json:"tenantId"`
}
