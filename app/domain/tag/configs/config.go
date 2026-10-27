package tagconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string `yaml:"mongoDbKey"`
	RedisDbKey string `yaml:"redisDbKey"`
	S3Key      string `yaml:"s3Key"`
	Enable     bool   `yaml:"enable"`
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

func GetRedisDbKey() string {
	return cfg.RedisDbKey
}
func GetS3Key() string {
	return cfg.S3Key
}
func GetEnable() bool {
	return cfg.Enable
}
