package cmdwrite

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type TempCreateCommand struct {
	CommandId   string                `json:"commandId"     validate:"gt=0"` // 命令ID
	IsValidOnly bool                  `json:"isValidOnly"   validate:"gt=0"` // 是否仅验证，不执行
	Data        field.TempCreateField `json:"data"`
}

// GetAggregateId
// @Description: 获取聚合根Id
func (c *TempCreateCommand) GetAggregateId() ddd.AggregateId {
	return ddd.NewAggregateId(c.Data.Id)
}

// GetCommandId
// @Description: 获取命令Id
func (c *TempCreateCommand) GetCommandId() string {
	return c.CommandId
}

// GetTenantId
// @Description: 获取租户Id
func (c *TempCreateCommand) GetTenantId() string {
	return c.Data.TenantId
}

// GetIsValidOnly
// @Description: 是否只验证不执行。
func (c *TempCreateCommand) GetIsValidOnly() bool {
	return c.IsValidOnly
}

// Validate
// @Description: 命令数据验证
func (c *TempCreateCommand) Validate() error {
	ve := ddd.ValidateCommand(c, nil)
	if len(c.Data.MasterType) == 0 {
		ve.AppendField("MasterType", "【数据类型】不能为空")
	}
	if len(c.Data.MapHeads) == 0 {
		ve.AppendField("MapHeads", "【表头】不能为空")
	}
	if len(c.Data.Fields) == 0 {
		ve.AppendField("Fields", "【转换规则】不能为空")
	}
	if len(c.Data.FileName) == 0 {
		ve.AppendField("FileName", "【文件名称】不能为空")
	}
	if len(c.Data.FileId) == 0 {
		ve.AppendField("FileId", "【文件ID】不能为空")
	}
	if len(c.Data.CaseId) == 0 {
		ve.AppendField("CaseId", "【案件ID】不能为空")
	}
	for i, head := range c.Data.MapHeads {
		for j, v := range head.Columns {
			if len(v.Label) == 0 {
				ve.AppendField(fmt.Sprintf("MapHead[%v].Columns[%v].Label", i, j), "【表头名称】不能为空")
			}
			if len(v.Key) == 0 {
				ve.AppendField(fmt.Sprintf("MapHead[%v].Columns[%v].Key", i, j), "【表头关键字】不能为空")
			}
		}
	}

	for i, v := range c.Data.Fields {
		if len(v.Key) == 0 {
			ve.AppendField(fmt.Sprintf("Fields[%v].Key", i), "【字段关键字】不能为空")
		}
		if len(v.Name) == 0 {
			ve.AppendField(fmt.Sprintf("Fields[%v].Name", i), "【字段名称】不能为空")
		}
	}
	return ve.GetError()
}
