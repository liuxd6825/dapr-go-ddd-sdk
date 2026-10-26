package graphconfigs

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

type Config struct {
	LivyServer LivyServer `yaml:"livyServer"`
	MongoDbKey string     `yaml:"mongoDbKey"`
	HugeDbKey  string     `yaml:"hugeDbKey"`
	DorisDbKey string     `yaml:"dorisDbKey"`
	Neo4jDbKey string     `yaml:"neo4jDbKey"`
	RedisDbKey string     `yaml:"redisDbKey"`
	S3Key      string     `yaml:"s3Key"`
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

func GetGraphDbKey() string {
	return cfg.Neo4jDbKey
}

func GetNeo4jDbKey() string {
	return cfg.Neo4jDbKey
}

func GetMongoDbKey() string {
	return cfg.MongoDbKey
}

func GetS3Key() string {
	return cfg.S3Key
}

func GetHugeDbKey() string {
	return cfg.HugeDbKey
}

func GetDorisDbKey() string {
	return cfg.DorisDbKey
}
