package drawconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string `yaml:"mongoDbKey"`
	HugeDbKey  string `yaml:"hugeDbKey"`
	DorisDbKey string `yaml:"dorisDbKey"`
	RedisDbKey string `yaml:"redisDbKey"`
}

var cfg Config
var once sync.Once

func GetConfig() Config {
	return cfg
}

func Init(env *envs.Env) {
	once.Do(func() {
		env.GetServiceConfig("draw-service", &cfg)
	})
}
