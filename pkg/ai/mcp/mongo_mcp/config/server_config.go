package config

type ServerConfig struct {
	AppName        string `json:"app_name"`
	Addr           string `json:"addr"`
	RsName         string `json:"rs_name"`
	UserName       string `json:"user_name"`
	Password       string `json:"password"`
	AutoDB         string `json:"auto_db"`
	AuthMechanism  string `json:"auth_mechanism"`
	ConnectTimeout int    `json:"connect_timeout"`
	DBName         string `json:"db_name"`

	LlmBinding     string `json:"llm_binding"`
	LlmModel       string `json:"llm_model"`
	LlmBindingHost string `json:"llm_binding_host"`
	LlmApiKey      string `json:"llm_api_key"`

	Tools []ToolConfig `json:"tools"`
}
