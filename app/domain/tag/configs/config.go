package tagconfigs

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"

type Config struct {
	MongoDbKey string `json:"mongoDbKey"`
	S3Key      string `json:"s3Key"`
	RedisKey   string `json:"redisKey"`
}

var cfg Config

func Init(env *envs.Env) {
	env.GetAppSubItem("sysService", &cfg)
}

func GetConfig() Config {
	return cfg
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}
