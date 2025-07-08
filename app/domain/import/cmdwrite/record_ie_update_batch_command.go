package cmdwrite

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

// RecordIeUpdateFilterCommand
// @Description:
type RecordIeUpdateFilterCommand struct {
	CommandId   string                           `json:"commandId"  validate:"required"` // 命令ID
	IsValidOnly bool                             `json:"isValidOnly"  validate:"-"`      // 是否仅验证，不执行
	Data        field.RecordIeUpdateFilterFields `json:"data"  validate:"required"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *RecordIeUpdateFilterCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId("")
}

// GetCommandId
// @Description: 获取命令Id
func (c *RecordIeUpdateFilterCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *RecordIeUpdateFilterCommand) GetTenantId() string {
	return c.Data.TenantId
}

// GetIsValidOnly
// @Description: 是否只验证不执行。
func (c *RecordIeUpdateFilterCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

// Validate
// @Description: 命令数据验证
func (c *RecordIeUpdateFilterCommand) Validate() error {
	ve := errors.NewVerifyError()
	if len(c.Data.TenantId) == 0 {
		ve.AppendField("tenantId", "不能为空")
	}
	if len(c.Data.Values) == 0 {
		ve.AppendField("values", "不能为空")
	}
	if len(c.Data.TaskId) == 0 {
		ve.AppendField("taskId", "不能为空")
	}
	return ve.GetError()
}
