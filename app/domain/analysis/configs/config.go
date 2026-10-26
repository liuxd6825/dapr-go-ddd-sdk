package analysisconfigs

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"

type Config struct {
	MongoDbKey string `yaml:"mongoDbKey"`
	DorisDbKey string `yaml:"dorisDbKey"`
	RedisDbKey string `yaml:"redisDbKey"`
	S3Key      string `yaml:"s3Key"`
}

var cfg Config

func Init(env *envs.Env) {
	env.GetServiceConfig("analysis-service", &cfg)
}

func GetConfig() Config {
	return cfg
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}

func GetDorisDbKey() string {
	return cfg.DorisDbKey
}
