package ragconfigs

import (
	"encoding/json"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string    `json:"mongoDbKey"`
	HugeDbKey  string    `json:"hugeDbKey"`
	Rag        RagConfig `json:"rag"`
}
type RagConfig struct {
	LLM      LLMConfig      `yaml:"llm" json:"llm"`
	Embedder EmbedderConfig `yaml:"embedder" json:"embedder"`
	Vector   VectorConfig   `yaml:"vector" json:"vector"`
}

type LLMConfig struct {
	Type    string `json:"type"` // ollama;
	Model   string `json:"model"`
	APIKey  string `json:"apiKey"`
	BaseUrl string `json:"baseUrl"`
}

type EmbedderConfig struct {
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
	ApiKey  string `json:"apiKey"`
}

type VectorConfig struct {
	Addr           string `json:"addr"`
	CollectionName string `json:"collectionName"`
	Dim            int    `json:"dim"`
}

var cfg Config

func Init(env *envs.Env) {
	env.GetAppSubItem("ragService", &cfg)
}

func GetConfig() Config {
	return cfg
}

func ReadRagConfig(ragMeta any) (*RagConfig, error) {
	jsonData, err := json.Marshal(ragMeta)
	if err != nil {
		return nil, err
	}
	var ragCfg RagConfig
	if err := json.Unmarshal(jsonData, &ragCfg); err != nil {
		return nil, err
	}
	return &ragCfg, nil
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}
