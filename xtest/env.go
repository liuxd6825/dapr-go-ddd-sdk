package xtest

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
)

type EnvConfig struct {
}

func NewEnvConfig() *env.Env {
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = "localhost"
	res.App.HttpPort = 0
	return res
}

func NewEnvConfig_Neo4j(ipAddr ...string) *env.Env {
	addr := "localhost"
	dbKey := "neo4j"
	for _, ip := range ipAddr {
		if ip != "" {
			addr = ip
		}
	}
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = "localhost"
	res.App.HttpPort = 0
	res.App.Meta = map[string]any{
		"db": dbKey,
	}

	res.AddNeo4j(&env.Neo4j{
		DbKey:    dbKey,
		Host:     addr,
		Port:     "7687",
		Database: "",
		UserName: "neo4j",
		Password: "12345678",
	})
	res.Init()
	return res
}

func InitEnv_MongoRemoteTest(opts ...*MongoOptions) *env.Env {
	envVal := NewEnvConfigMongo(GetMongoEnv_Remote(opts...))
	env.SetEnv(envVal)
	return envVal
}

func InitEnv_MongoRemoteMaster(opts ...*MongoOptions) *env.Env {
	opt := NewMongoOptions(opts...).SetDBName("master")
	envVal := NewEnvConfigMongo(GetMongoEnv_Remote(opt))
	env.SetEnv(envVal)
	return envVal
}

func InitEnv_MongoLocal(opts ...*MongoOptions) *env.Env {
	envVal := NewEnvConfigMongo(GetMongoEnv_Local(opts...))
	env.SetEnv(envVal)
	return envVal
}

func NewEnvConfigMongo(mongoCfg *env.Mongo) *env.Env {
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = "localhost"
	res.App.HttpPort = 0
	res.App.Meta = map[string]any{
		"db": mongoCfg.DbKey,
	}
	docFs := map[string]any{
		"name": "documentIoStore",
		"type": "local",
		"path": "/Users/lxd/Projects/duxm/h-master/document",
	}
	res.Fs = append(res.Fs, docFs)

	res.AddMongo(mongoCfg)
	res.Init()
	return res
}

func NewEnvConfig_Fs(name, path string) *env.Env {
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = "localhost"
	res.App.HttpPort = 0
	res.Fs = []map[string]any{
		{
			"name": name,
			"type": "local",
			"path": path,
		},
	}
	res.Init()
	return res
}

func (e EnvConfig) GetAppId() string {
	return "test"
}

func (e EnvConfig) GetAppName() string {
	return "app"
}

func (e EnvConfig) GetAppHttpHost() string {
	return "localhost"
}

func (e EnvConfig) GetAppHttpPort() int {
	return 0
}

func (e EnvConfig) GetDaprHost() string {
	return "localhost"
}

func (e EnvConfig) GetDaprHttpPort() int64 {
	return 0
}

func (e EnvConfig) GetDaprGrpcPort() int64 {
	return 0
}

func (e EnvConfig) GetFsManager() *fsm.Manager {
	return nil
}

func (e EnvConfig) GetHServerSrcPath() string {
	return ""
}

func (e EnvConfig) GetHServerEnable() bool {
	return false
}
