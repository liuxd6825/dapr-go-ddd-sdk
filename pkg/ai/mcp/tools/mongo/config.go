package mongo

import "time"

type Config struct {
	Kind string     `json:"kind" yaml:"kind"`
	DB   DBConfig   `json:"db" yaml:"db" `
	LLM  LLMConfig  `json:"llm" yaml:"llm"`
	Tool ToolConfig `json:"tool" yaml:"tool"`
}

type DBConfig struct {
	Host                   string        `json:"host" yaml:"host"`
	ReplicaSet             string        `json:"replicaSet" yaml:"replicaSet"`
	UserName               string        `json:"userName" yaml:"userName"`
	Password               string        `json:"password" yaml:"password"`
	AuthSource             string        `json:"authSource" yaml:"authSource"`
	AuthMechanism          string        `json:"authMechanism" yaml:"authMechanism"`
	Options                string        `json:"options" yaml:"options"`
	DBName                 string        `json:"dbName" yaml:"dbName"`
	CollName               string        `json:"collName" yaml:"collName"`
	MaxPoolSize            uint64        `json:"maxPoolSize" yaml:"maxPoolSize"`
	ConnectTimeout         time.Duration `json:"connectTimeout" yaml:"connectTimeout"`
	OperationTimeout       uint          `json:"operationTimeout" yaml:"operationTimeout"`
	ServerSelectionTimeout time.Duration `json:"serverSelectionTimeout" yaml:"serverSelectionTimeout"`
	SocketTimeout          time.Duration `json:"socketTimeout" yaml:"socketTimeout"`
	HeartbeatInterval      time.Duration `json:"heartbeatInterval" yaml:"heartbeatInterval"`
	LocalThreshold         time.Duration `json:"localThreshold" yaml:"localThreshold"`
	MaxConnIdleTime        time.Duration `json:"maxConnIdleTime" yaml:"maxConnIdleTime"`
	Direct                 *bool         `json:"direct" yaml:"direct"`
}

type LLMType string

const (
	LLMType_Ollama LLMType = "ollama"
	LLMType_OpenAi LLMType = "openAi"
)

type LLMConfig struct {
	Type    LLMType `json:"type" yaml:"type"` // ollama
	BaseURL string  `json:"baseURL" yaml:"baseURL"`
	Model   string  `json:"model" yaml:"model"`
	ApiKey  string  `json:"api_key" yaml:"apiKey"`
}

type ToolConfig struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Prompt      string `json:"prompt" yaml:"prompt"`
	Schema      string `json:"schema" yaml:"schema"`
}
