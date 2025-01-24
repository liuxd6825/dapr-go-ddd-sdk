package element

type RequestConfig struct {
	Type          string        `json:"type"`
	Name          string        `json:"name"`
	URL           string        `json:"url"`
	AbsURl        string        `json:"absUri"`
	Description   string        `json:"description"`
	ParamsType    string        `json:"paramsType"`
	LinkParamsUrl string        `json:"linkParamsUrl"`
	Script        *ScriptConfig `json:"script"`
}

type Request interface {
	GetParamsValue(wctx WebContext) map[string]any
	Config() *RequestConfig
	Initialize() error
}
