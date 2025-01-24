package request

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/schema"
)

type Handle func(cxt *element.WebContext, data any)
type MethodType string

const (
	GET    MethodType = "GET"
	POST   MethodType = "POST"
	PUT    MethodType = "PUT"
	DELETE MethodType = "DELETE"
)

func (m MethodType) String() string {
	return string(m)
}

type InParamType string

func (f InParamType) String() string {
	return string(f)
}

const (
	InParamTypeURL        InParamType = "url"        // URL中的参数
	InParamTypePath       InParamType = "path"       // 路径参数
	InParamTypeBody       InParamType = "body"       // 请求body参数
	InParamTypeFormValue  InParamType = "formValue"  // 从FormData中读取string
	InParamTypeFormObject InParamType = "formObject" // 从FormData中读取json转成对象
	InParamTypeFormFile   InParamType = "formFile"   // 从FormFile中读取文件
)

type RequestData struct {
	Data   *goja.Object   `json:"runValues"` // 从body中读取的map数据
	Params map[string]any `json:"params"`    // 参数
}

type HandleOptions struct {
	Method      MethodType                                               `json:"method"`      // 请求类型
	Path        string                                                   `json:"path"`        // 请求的URL
	Description string                                                   `json:"description"` // 方法说明
	Body        *schema.Schema                                           `json:"body"`        // 请求时body的数据定义
	Params      map[string]common.RequestParam                           `json:"params"`      // 参数定义
	ParamsUrl   string                                                   `json:"paramsUrl"`   // 从URL中加载params的定义
	HandleName  string                                                   `json:"handleName"`  // 名称
	Handle      func(cxt *element.WebContext, params map[string]any) any `json:"-"`           // 控制器
}

func (h HandleOptions) GetMethod() MethodType {
	return h.Method
}

func (h *HandleOptions) GetPath() string {
	return h.Path
}

func (h *HandleOptions) GetDescription() string {
	return h.Description
}

func (h *HandleOptions) GetBody() *schema.Schema {
	return h.Body
}

func (h *HandleOptions) GetParams() map[string]common.RequestParam {
	return h.Params
}
