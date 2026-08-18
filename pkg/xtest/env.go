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

type Neo4jOptions struct {
	Addr     string
	DBKey    string
	Port     string
	Database string
	UserName string
	Password string
}

func NewEnvConfigNeo4j(opts ...Neo4jOptions) *env.Env {
	addr := "localhost"
	dbKey := "neo4j"
	port := "7687"
	database := ""
	userName := "neo4j"
	password := "12345678"
	for _, opt := range opts {
		if opt.Addr != "" {
			addr = opt.Addr
		}
		if opt.DBKey != "" {
			dbKey = opt.DBKey
		}
		if opt.Port != "" {
			port = opt.Port
		}
		if opt.Database != "" {
			database = opt.Database
		}
		if opt.UserName != "" {
			userName = opt.UserName
		}
		if opt.Password != "" {
			password = opt.Password
		}
	}
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = ""
	res.App.HttpPort = 0
	res.App.Meta = map[string]any{
		"db": dbKey,
	}

	res.AddNeo4j(&env.Neo4j{
		DbKey:    dbKey,
		Host:     addr,
		Port:     port,
		Database: database,
		UserName: userName,
		Password: password,
	})
	res.Init()
	return res
}

var Neo4jRemoveOption = Neo4jOptions{
	Addr:     "192.168.120.224",
	Port:     "7687",
	Database: "",
	Password: "12345678",
	UserName: "neo4j",
}

var Neo4jLocalOption = Neo4jOptions{
	Addr:     "127.0.0.1",
	Port:     "7687",
	Database: "",
	Password: "12345678",
	UserName: "neo4j",
}

func InitEnv_Neo4j(opts ...Neo4jOptions) *env.Env {
	envVal := NewEnvConfigNeo4j(opts...)
	env.SetEnv(envVal)
	return envVal
}

type NebulaOptions struct {
	Addr     string
	DBKey    string
	User     string
	Password string
	Space    string
	PoolSize int
}

var NebulaRemoveOption = NebulaOptions{
	Addr:     NebulaHostRemote + ":9669",
	DBKey:    NebulaDBKey,
	User:     "root",
	Password: "nebula",
	Space:    "test_default",
	PoolSize: 10,
}

var NebulaLocalOption = NebulaOptions{
	Addr:     NebulaHostLocal + ":9669",
	DBKey:    NebulaDBKey,
	User:     "root",
	Password: "nebula",
	Space:    "test_default",
	PoolSize: 10,
}

func NewEnvConfigNebula(opts ...NebulaOptions) *env.Env {
	addr := NebulaHostLocal + ":9669"
	dbKey := NebulaDBKey
	user := "root"
	password := "nebula"
	space := "test_default"
	poolSize := 10
	for _, opt := range opts {
		if opt.Addr != "" {
			addr = opt.Addr
		}
		if opt.DBKey != "" {
			dbKey = opt.DBKey
		}
		if opt.User != "" {
			user = opt.User
		}
		if opt.Password != "" {
			password = opt.Password
		}
		if opt.Space != "" {
			space = opt.Space
		}
		if opt.PoolSize > 0 {
			poolSize = opt.PoolSize
		}
	}
	res := env.NewEnv()
	res.App.AppId = "test"
	res.App.AppName = "app"
	res.App.HttpHost = ""
	res.App.HttpPort = 0
	res.App.Meta = map[string]any{
		"db": dbKey,
	}

	res.AddNebula(&env.Nebula{
		Name:     dbKey,
		Addrs:    []string{addr},
		User:     user,
		Password: password,
		Space:    space,
		PoolSize: poolSize,
	})
	res.Init()
	return res
}

func InitEnv_Nebula(opts ...NebulaOptions) *env.Env {
	envVal := NewEnvConfigNebula(opts...)
	env.SetEnv(envVal)
	return envVal
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

func init() {
	InitTimeZone()
}
