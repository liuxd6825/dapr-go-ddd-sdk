package sysconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string    `yaml:"mongoDbKey"`
	S3Key      string    `yaml:"s3Key"`
	RedisDbKey string    `yaml:"redisDbKey"`
	Org        OryConfig `yaml:"org"`
}

type OryConfig struct {
	Url string `yaml:"url"`
}

var cfg Config
var once sync.Once

func Init(env *envs.Env) {
	once.Do(func() {
		env.GetServiceConfig("sys-service", &cfg)
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
