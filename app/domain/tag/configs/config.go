package tagconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string `yaml:"mongoDbKey"`
	S3Key      string `yaml:"s3Key"`
	RedisDbKey string `yaml:"redisDbKey"`
}

var cfg Config
var once sync.Once

func Init(env *envs.Env) {
	once.Do(func() {
		env.GetServiceConfig("tag-service", &cfg)
	})

}

func GetConfig() Config {
	return cfg
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}
