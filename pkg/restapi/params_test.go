package restapi

import (
	"testing"
)

type Base[T any] struct {
	CommandId   string   `json:"commandId"  validate:"required" ` // 命令ID
	IsValidOnly bool     `json:"isValidOnly"  validate:"-"`       // 是否仅验证，不执行
	UpdateMask  []string `json:"updateMask"  validate:"-"`
	Data        T        `json:"data"  validate:"required"`
}

type Command struct {
	Base[TestCommandData]
}
type TestCommandData struct {
	Name string `json:"name" validate:"required" title:"姓名"`
}

func Test_Validate(t *testing.T) {

}
