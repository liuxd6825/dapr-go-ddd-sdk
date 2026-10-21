package importconfigs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string     `json:"mongoDbKey"`
	DorisDbKey string     `json:"dorisDbKey"`
	LivyServer LivyServer `json:"livyServer"`
	Rag        RagConfig  `json:"rag"`
}
type RagConfig struct {
	LLM      LLMConfig      `json:"llm"`
	Embedder EmbedderConfig `json:"embedder"`
	Vector   VectorConfig   `json:"vector"`
}

type LivyServer struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type LLMConfig struct {
	Type      string `json:"type"` // ollama;
	Model     string `json:"model"`
	APIKey    string `json:"apiKey"`
	BaseUrl   string `json:"baseUrl"`
	CtxLength int    `json:"ctxLength"`
	Think     string `json:"think"`
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

var _cfg Config

func GetConfig() Config {
	return _cfg
}

func Init(env *envs.Env) {
	env.GetAppSubItem("importService", &_cfg)
}
