package ragconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string         `yaml:"mongoDbKey"`
	HugeDbKey  string         `yaml:"hugeDbKey"`
	Neo4jDbKey string         `yaml:"neo4jDbKey"`
	RedisDbKey string         `yaml:"redisDbKey"`
	Rag        RagConfig      `yaml:"rag"`
	LightRag   LightRagConfig `yaml:"lightRag"`
	Enable     bool           `yaml:"enable"`
}
type RagConfig struct {
	LLM      LLMConfig      `yaml:"llm" yaml:"llm"`
	Embedder EmbedderConfig `yaml:"embedder" yaml:"embedder"`
	Vector   VectorConfig   `yaml:"vector" yaml:"vector"`
}

type LightRagConfig struct {
	Url    string `yaml:"url"`
	ApiKey string `yaml:"apiKey"`
}

type LLMConfig struct {
	Type    string `yaml:"type"` // ollama;
	Model   string `yaml:"model"`
	APIKey  string `yaml:"apiKey"`
	BaseUrl string `yaml:"baseUrl"`
}

type EmbedderConfig struct {
	BaseURL string `yaml:"baseUrl"`
	Model   string `yaml:"model"`
	ApiKey  string `yaml:"apiKey"`
}

type VectorConfig struct {
	Addr           string `yaml:"addr"`
	CollectionName string `yaml:"collectionName"`
	Dim            int    `yaml:"dim"`
}

var cfg Config
var once sync.Once

func Init(env *envs.Env) {
	once.Do(func() {
		env.GetServiceConfig("rag-service", &cfg)
	})
}

func GetConfig() Config {
	return cfg
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}

func GetGraphDbKey() string {
	return cfg.Neo4jDbKey
}

func GetEnable() bool {
	return cfg.Enable
}
