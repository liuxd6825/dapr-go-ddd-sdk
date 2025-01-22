package hserver

type RequestConfig struct {
	Type          string       `json:"type"`
	Name          string       `json:"name"`
	URL           string       `json:"url"`
	AbsURl        string       `json:"abs_uri"`
	Description   string       `json:"description"`
	ParamsType    string       `json:"params_type"`
	LinkParamsUrl string       `json:"link_params_url"`
	Script        ScriptConfig `json:"script"`
}
