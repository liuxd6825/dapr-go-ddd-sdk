package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type TempDeleteCommand struct {
	CommandId   string                `json:"commandId"     validate:"gt=0"` // 命令ID
	IsValidOnly bool                  `json:"isValidOnly"   validate:"gt=0"` // 是否仅验证，不执行
	Data        field.TempDeleteField `json:"data"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *TempDeleteCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(c.Data.Id)
}

// GetCommandId
// @Description: 获取命令Id
func (c *TempDeleteCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *TempDeleteCommand) GetTenantId() string {
	return c.Data.TenantId
}

// GetIsValidOnly
// @Description: 是否只验证不执行。
func (c *TempDeleteCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

// Validate
// @Description: 命令数据验证
func (c *TempDeleteCommand) Validate() error {
	ve := ddd.ValidateCommand(c, nil)
	if len(c.Data.Id) == 0 {
		ve.AppendField("Data.Id", "不能为空")
	}
	if len(c.Data.TenantId) == 0 {
		ve.AppendField("Data.TenantId", "不能为空")
	}
	return ve.GetError()
}
