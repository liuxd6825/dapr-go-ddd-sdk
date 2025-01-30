package element

type SubEventConfig struct {
	Name          string        `json:"name"`
	ParamsType    string        `json:"paramsType"`
	Description   string        `json:"description"`
	LinkParamsUrl string        `json:"linkParamsUrl"`
	Script        *ScriptConfig `json:"script"`
	Version       string        `json:"version"`
}

type SubEvent interface {
	Config() *SubEventConfig
	Initialize() error
}

type SubServiceConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	AppId       string `json:"appId"`
	Pubsub      string `json:"pubsub"`
}

type SubService interface {
	Config() *SubServiceConfig
	Initialize() error
}
