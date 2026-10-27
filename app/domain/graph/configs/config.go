package graphconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	LivyServer LivyServer `yaml:"livyServer"`
	MongoDbKey string     `yaml:"mongoDbKey"`
	DorisDbKey string     `yaml:"dorisDbKey"`
	GraphDbKey string     `yaml:"graphDbKey"`
	ImdbKey    string     `yaml:"imdbKey"`
	OssKey     string     `yaml:"ossKey"`
	Enable     bool       `yaml:"enable"`
}

type LivyServer struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

var cfg Config
var once sync.Once

func GetConfig() Config {
	return cfg
}

func Init(env *envs.Env) {
	once.Do(func() {
		env.GetServiceConfig("graph-service", &cfg)
	})

}

func GetEnable() bool {
	return cfg.Enable
}

func GetGraphDbKey() string {
	return cfg.GraphDbKey
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}

func GetOssKey() string {
	return cfg.OssKey
}

func GetDorisDbKey() string {
	return cfg.DorisDbKey
}
