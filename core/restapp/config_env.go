package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mapperutils"
)

func NewEnv(cfg *EnvConfig) *env.Env {
	env := env.NewEnv()
	env.Fs = cfg.Fs
	env.Name = cfg.Name
	env.App = NewApp(cfg.App)
	env.Fsm = NewFsm(cfg.Fs, cfg.App.HServer.SrcName)
	env.Log = NewLog(cfg.Log)
	env.Dapr = NewDapr(cfg.Dapr)
	env.Auth = NewAuth(cfg.Auth)
	env.Mysql = NewMySQL(cfg.Mysql)
	env.Mongo = NewMongo(cfg.Mongo)
	env.Neo4j = NewNeo4j(cfg.Neo4j)
	env.Minio = NewMinio(cfg.Minio)
	env.Redis = NewRedis(cfg.Redis)
	return env
}

func NewRedis(cfg map[string]*RedisConfig) map[string]*env.Redis {
	items := map[string]*env.Redis{}
	for key, item := range cfg {
		m := &env.Redis{}
		copy(item, m)
		items[key] = m
	}
	return items
}

func NewMinio(cfg map[string]*MinioConfig) map[string]*env.Minio {
	items := map[string]*env.Minio{}
	for key, item := range cfg {
		m := env.NewMinio()
		copy(item, m)
		items[key] = m
	}
	return items
}

func NewNeo4j(cfg map[string]*Neo4jConfig) map[string]*env.Neo4j {
	items := map[string]*env.Neo4j{}
	for key, item := range cfg {
		neo4j := env.NewNeo4j()
		copy(item, neo4j)
		items[key] = neo4j
	}
	return items
}

func NewMySQL(cfg map[string]*MySqlConfig) map[string]*env.MySql {
	items := map[string]*env.MySql{}
	for key, item := range cfg {
		mysql := env.NewMySQL()
		copy(item, mysql)
		items[key] = mysql
	}
	return items
}

func NewMongo(cfg map[string]*MongoConfig) map[string]*env.Mongo {
	items := map[string]*env.Mongo{}
	for key, item := range cfg {
		mysql := &env.Mongo{}
		copy(item, mysql)
		items[key] = mysql
	}
	return items
}

func NewDapr(cfg *DaprConfig) *env.Dapr {
	dapr := env.NewDapr()
	copy(cfg, dapr)
	return dapr
}

func NewAuth(cfg *AuthConfig) *env.Auth {
	auth := env.NewAuth()
	copy(cfg, auth)
	return auth
}

func NewApp(cfg *AppConfig) *env.App {
	app := env.NewApp()
	copy(cfg, app)
	return app
}

func NewLog(cfg *LogConfig) *env.Log {
	log := env.NewLog()
	copy(cfg, log)
	return log
}

func NewFsm(fs []map[string]any, defFsName string) *fsm.Manager {
	if len(fs) != 0 {
		fsManager, err := fsm.NewManagerWithConfigs(fs, defFsName)
		if err != nil {
			panic(errors.New("fs.NewManagerWithConfigs() err: %s", err.Error()))
		}
		return fsManager
	}
	return fsm.NewManager()
}

func copy(fromObj any, toObj any) {
	err := mapperutils.Mapper(fromObj, toObj)
	if err != nil {
		panic(fmt.Sprintf("NewApp() error:%s", err.Error()))
	}
}
