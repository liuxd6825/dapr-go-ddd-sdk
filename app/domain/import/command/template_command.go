package command

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type TempCreateCommand struct {
	xbase.Command[field.TempCreateField]
}
type TempUpdateCommand struct {
	xbase.Command[field.TempUpdateField]
}
type TempDeleteCommand struct {
	xbase.Command[field.TempDeleteField]
}

// Validate
// @Description: 命令数据验证
func (c *TempCreateCommand) Validate() error {
	ve := errors.NewVerifyError()
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

// Validate
// @Description: 命令数据验证
func (c *TempUpdateCommand) Validate() error {
	ve := errors.NewVerifyError()
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
