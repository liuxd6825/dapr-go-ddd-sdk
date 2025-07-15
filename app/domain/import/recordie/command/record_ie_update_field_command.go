package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

// RecordIeUpdateFieldCommand
// @Description:
type RecordIeUpdateFieldCommand struct {
	CommandId   string                          `json:"commandId"  validate:"required"` // 命令ID
	IsValidOnly bool                            `json:"isValidOnly"  validate:"-"`      // 是否仅验证，不执行
	Data        field.RecordIeUpdateFieldFields `json:"data"  validate:"required"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *RecordIeUpdateFieldCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(c.Data.Id)
}

// GetCommandId
// @Description: 获取命令Id
func (c *RecordIeUpdateFieldCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *RecordIeUpdateFieldCommand) GetTenantId() string {
	return c.Data.TenantId
}

// GetIsValidOnly
// @Description: 是否只验证不执行。
func (c *RecordIeUpdateFieldCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

// Validate
// @Description: 命令数据验证
func (c *RecordIeUpdateFieldCommand) Validate() error {
	ve := ddd.ValidateCommand(c, nil)
	if len(c.Data.Id) == 0 {
		ve.AppendField("id", "不能为空")
	}
	if len(c.Data.TenantId) == 0 {
		ve.AppendField("tenantId", "不能为空")
	}
	if len(c.Data.Values) == 0 {
		ve.AppendField("values", "不能为空")
	}
	return ve.GetError()
}
