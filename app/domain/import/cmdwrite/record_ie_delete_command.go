package cmdwrite

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

// RecordIeDeleteCommand
// @Description:
type RecordIeDeleteCommand struct {
	CommandId   string                     `json:"commandId"  validate:"required"` // 命令ID
	IsValidOnly bool                       `json:"isValidOnly"`
	Data        field.RecordIeDeleteFields `json:"data"  validate:"required"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *RecordIeDeleteCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(c.Data.TaskId)
}

// GetIsValidOnly
// @Description:
// @receiver c
// @return bool
func (c *RecordIeDeleteCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

// GetCommandId
// @Description: 获取命令Id
func (c *RecordIeDeleteCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *RecordIeDeleteCommand) GetTenantId() string {
	return c.Data.TenantId
}

func (c *RecordIeDeleteCommand) IsAggregateCreateCommand() {

}

// Validate
// @Description: 命令数据验证
func (c *RecordIeDeleteCommand) Validate() error {
	ve := ddd.ValidateCommand(c, nil)
	if len(c.GetTenantId()) == 0 {
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
