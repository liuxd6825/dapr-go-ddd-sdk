package ddd_command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
)

type DeleteByIdCommand struct {
	CommandId   string                `json:"commandId"  validate:"required"` // 命令ID
	IsValidOnly bool                  `json:"isValidOnly"  validate:"-"`      // 是否仅验证，不执行
	Data        DeleteByIdCommandData `json:"data"`
}

type DeleteByIdCommandData struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
}

func NewDeleteByIdCommand() *DeleteByIdCommand {
	return &DeleteByIdCommand{}
}

func (d *DeleteByIdCommand) GetCommandId() string {
	return d.CommandId
}

func (d *DeleteByIdCommand) GetTenantId() string {
	return d.Data.TenantId
}

func (d *DeleteByIdCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(d.Data.Id)
}

func (d *DeleteByIdCommand) GetIsValidOnly() bool {
	return d.IsValidOnly
}

func (d *DeleteByIdCommand) Validate() error {
	verify := errors.NewVerifyError()
	if len(d.CommandId) == 0 {
		verify.AppendField("CommandId", "不能为空")
	}
	if len(d.Data.TenantId) == 0 {
		verify.AppendField("TenantId", "不能为空")
	}
	if len(d.Data.TenantId) == 0 {
		verify.AppendField("ID", "不能为空")
	}
	if len(verify.Errors) > 0 {
		return verify
	}
	return nil
}
