package importconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string     `yaml:"mongoDbKey"`
	DorisDbKey string     `yaml:"dorisDbKey"`
	RedisDbKey string     `yaml:"redisDbKey"`
	LivyServer LivyServer `yaml:"livyServer"`
	Rag        RagConfig  `yaml:"rag"`
}
type RagConfig struct {
	LLM      LLMConfig      `yaml:"llm"`
	Embedder EmbedderConfig `yaml:"embedder"`
	Vector   VectorConfig   `yaml:"vector"`
}

type LivyServer struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type LLMConfig struct {
	Type      string `yaml:"type"` // ollama;
	Model     string `yaml:"model"`
	APIKey    string `yaml:"apiKey"`
	BaseUrl   string `yaml:"baseUrl"`
	CtxLength int    `yaml:"ctxLength"`
	Think     string `yaml:"think"`
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

func GetConfig() Config {
	return cfg
}

func Init(env *envs.Env) {
	once.Do(func() {
		env.GetServiceConfig("import-service", &cfg)
	})
}
