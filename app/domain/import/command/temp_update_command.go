package command

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type TempUpdateCommand struct {
	CommandId   string                `json:"commandId"     validate:"gt=0"` // 命令ID
	IsValidOnly bool                  `json:"isValidOnly"   validate:"gt=0"` // 是否仅验证，不执行
	UpdateMask  []string              `json:"updateMask"  validate:"-"`      // 要更新的字段项，空值：更新所有字段
	Data        field.TempCreateField `json:"data"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *TempUpdateCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(c.Data.Id)
}

// GetCommandId
// @Description: 获取命令Id
func (c *TempUpdateCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *TempUpdateCommand) GetTenantId() string {
	return c.Data.TenantId
}

// GetIsValidOnly
// @Description: 是否只验证不执行。
func (c *TempUpdateCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

// Validate
// @Description: 命令数据验证
func (c *TempUpdateCommand) Validate() error {
	ve := ddd.ValidateCommand(c, nil)
	if len(c.Data.Id) == 0 {
		ve.AppendField("Data.Id", "不能为空")
	}
	if len(c.Data.CaseId) == 0 {
		ve.AppendField("Data.CaseId", "不能为空")
	}
	if len(c.Data.TenantId) == 0 {
		ve.AppendField("Data.TenantId", "不能为空")
	}
	if len(c.Data.BankName) == 0 {
		ve.AppendField("Data.BankName", "不能为空")
	}
	if len(c.Data.MasterType) == 0 {
		ve.AppendField("Data.MasterType", "不能为空")
	}
	for i, head := range c.Data.MapHeads {
		for j, v := range head.Columns {
			if len(v.Label) == 0 {
				ve.AppendField(fmt.Sprintf("Data.MapHeads[%v].Columns[%v].Label", i, j), "不能为空")
			}
			if len(v.Key) == 0 {
				ve.AppendField(fmt.Sprintf("Data.MapHeads[%v].Columns[%v].Key", i, j), "不能为空")
			}
		}
	}
	for i, v := range c.Data.Fields {
		if len(v.Key) == 0 {
			ve.AppendField(fmt.Sprintf("Data.Fields[%v].Key", i), "不能为空")
		}
		if len(v.Name) == 0 {
			ve.AppendField(fmt.Sprintf("Data.Fields[%v].Name", i), "不能为空")
		}
	}
	return ve.GetError()
}
