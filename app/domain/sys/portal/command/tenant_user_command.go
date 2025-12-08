package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/portal/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

type TenantUserCreateCommand struct {
	xbase.Command[TenantUserCreateCommandData]
}

type TenantUserUpdateCommand struct {
	xbase.Command[model.TenantUser]
}

type TenantUserDeleteBatchCommand struct {
	xbase.Command[DeleteByIds]
}

type DeleteByIds struct {
	Ids    []string `json:"ids"`
	CaseId string   `json:"caseId"`
}

type TenantUserCreateCommandData struct {
	TenIds  []string `json:"tenIds" required:"true"`
	UserIds []string `json:"userIds" required:"true"`
}
