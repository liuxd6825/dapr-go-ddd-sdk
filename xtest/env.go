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

func NewEnvConfig_Neo4j() *env.Env {
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = "localhost"
	res.App.HttpPort = 0

	res.AddNeo4j(&env.Neo4j{
		DbKey:    "neo4j",
		Host:     "localhost",
		Port:     "7687",
		Database: "",
		UserName: "neo4j",
		Password: "12345678",
	})
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
