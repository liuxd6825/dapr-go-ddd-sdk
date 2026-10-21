package docconfigs

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"

type Config struct {
	MongoDbKey string `json:"mongoDbKey"`
	S3Key      string `json:"s3Key"`
}

var cfg Config

func Init(env *envs.Env) {
	env.GetAppSubItem("analysisService", &cfg)
}

func GetConfig() Config {
	return cfg
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}
