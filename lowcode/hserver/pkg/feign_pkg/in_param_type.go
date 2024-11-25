package feign_pkg

import "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"

type InParamType string

const (
	InParamTypeURL  InParamType = "url"  // URL中的参数
	InParamTypePath InParamType = "path" // 路径参数
	InParamTypeBody InParamType = "body" // 请求body参数
)

func (f InParamType) String() string {
	return string(f)
}

type Param struct {
	In       string         `json:"in"` // InParamType
	Required bool           `json:"required"`
	Type     string         `json:"type"`
	Desc     string         `json:"desc"`
	Schema   *schema.Schema `json:"schema"`
}

type Options struct {
	Method string //MethodType
	Url    string
	Params map[string]Param
}
