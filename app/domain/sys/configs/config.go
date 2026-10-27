package sysconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	MongoDbKey string       `yaml:"mongoDbKey"`
	OssKey     string       `yaml:"ossKey"`
	ImdbKey    string       `yaml:"imdbKey"`
	Kratos     KratosConfig `yaml:"kratos"`
	Enable     bool         `yaml:"enable"`
}

type KratosConfig struct {
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

func GetImdbKey() string {
	return cfg.ImdbKey
}

func GetOssKey() string {
	return cfg.OssKey
}

func GetEnable() bool {
	return cfg.Enable
}
