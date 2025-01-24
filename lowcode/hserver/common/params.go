package common

import "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/xschema"

type ParamsType map[string]*RequestParam

func NewParamsType() *ParamsType {
	return &ParamsType{}
}

type RequestParam struct {
	In          string          `json:"in"`          // InParamType
	Required    bool            `json:"required"`    // 是否必填
	Type        string          `json:"type"`        // 数据类型
	Description string          `json:"description"` // 说明
	Example     any             `json:"example"`     // 数据示例
	Schema      *xschema.Schema `json:"schema"`      // 数据定义
	Default     any             `json:"default"`     // 默认值
}
