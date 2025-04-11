package schema

type HParamType string

const (
	HParamType_Path       HParamType = "path"       // 路径参数
	HParamType_Query      HParamType = "query"      // URL参数
	HParamTypee_Header    HParamType = "header"     // HTML head 参数
	HParamType_Body       HParamType = "body"       // 请求body参数
	HParamTypee_FormValue HParamType = "formValue"  // 从FormData中读取string
	HParamType_FormObject HParamType = "formObject" // 从FormData中读取json转成对象
	HParamType_FormFile   HParamType = "formFile"   // 从FormFile中读取文件
)

type HttpParam struct {
	Name    string     `json:"name"`
	Type    HParamType `json:"type"`
	Convert string     `json:"convert"`
}

func NewHttpParam() *HttpParam {
	return &HttpParam{}
}

func (p *HttpParam) init(values map[string]any) error {
	for k, v := range values {
		switch k {
		case "name":
			p.Name = v.(string)
		case "convert":
			p.Convert = v.(string)
		case "type":
			if str, ok := v.(string); ok {
				p.Type = HParamType(str)
			}
		}
	}
	return nil
}
