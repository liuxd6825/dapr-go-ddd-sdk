package graphconfigs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	LivyServer LivyServer `json:"livyServer"`
	MongoDbKey string     `json:"mongoDbKey"`
	S3Key      string     `json:"s3Key"`
	HugeDbKey  string     `json:"hugeDbKey"`
	DorisDbKey string     `json:"dorisDbKey"`
}

type LivyServer struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var _cfg Config

func GetConfig() Config {
	return _cfg
}

func Init(env *envs.Env) {
	env.GetAppSubItem("graphService", &_cfg)
}
