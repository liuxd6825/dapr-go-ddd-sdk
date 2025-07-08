package cmdwrite

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

// RecordIeUpdateCommand
// @Description:
type RecordIeUpdateCommand struct {
	CommandId   string                     `json:"commandId"  validate:"required"` // 命令ID
	IsValidOnly bool                       `json:"isValidOnly"   validate:"gt=0"`  // 是否仅验证，不执行
	UpdateMask  []string                   `json:"updateMask"`                     // 要更新的字段项，空值：更新所有字段
	Data        field.RecordIeUpdateFields `json:"data"  validate:"required"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *RecordIeUpdateCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(c.Data.TaskId)
}

// GetCommandId
// @Description: 获取命令Id
func (c *RecordIeUpdateCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *RecordIeUpdateCommand) GetTenantId() string {
	return c.Data.TenantId
}

// GetIsValidOnly
// @Description: 是否只验证不执行。
func (c *RecordIeUpdateCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

func (c *RecordIeUpdateCommand) IsAggregateCreateCommand() {

}

// Validate
// @Description: 命令数据验证
func (c *RecordIeUpdateCommand) Validate() error {
	ve := ddd.ValidateCommand(c, nil)
	if len(c.UpdateMask) == 0 {
		ve.AppendField("updateMask", "不能为空")
	}
	if len(c.Data.Id) == 0 {
		ve.AppendField("id", "不能为空")
	}
	if len(c.Data.TaskId) == 0 {
		ve.AppendField("taskId", "不能为空")
	}
	return ve.GetError()
}
