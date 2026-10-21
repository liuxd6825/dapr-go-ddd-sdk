package restapp

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/mapperutils"
)

func NewEnv(cfg *EnvConfig) *envs.Env {
	env := envs.NewEnv()
	env.Fs = cfg.Fs
	env.Name = cfg.Name

	env.App = NewApp(cfg.App)
	env.Fsm = NewFsm(cfg.Fs, cfg.App.HServer.SrcName)
	env.Log = NewLog(cfg.Log)
	env.Dapr = NewDapr(cfg.Dapr)
	//env.Auth = NewAuth(cfg.Auth)
	env.Mysql = NewMySQL(cfg.Mysql)
	env.Mongo = NewMongo(cfg.Mongo)
	env.Neo4j = NewNeo4j(cfg.Neo4j)
	env.Minio = NewMinio(cfg.Minio)
	env.Redis = NewRedis(cfg.Redis)
	env.Elastic = NewElastic(cfg.Elastic)
	return env
}

func NewElastic(cfg map[string]*ElasticConfig) map[string]*envs.Elastic {
	items := map[string]*envs.Elastic{}
	for key, item := range cfg {
		m := &envs.Elastic{}
		copy(item, m)
		items[key] = m
	}
	return items
}

func NewRedis(cfg map[string]*RedisConfig) map[string]*envs.Redis {
	items := map[string]*envs.Redis{}
	for key, item := range cfg {
		m := &envs.Redis{}
		copy(item, m)
		items[key] = m
	}
	return items
}

func NewMinio(cfg map[string]*MinioConfig) map[string]*envs.Minio {
	items := map[string]*envs.Minio{}
	for key, item := range cfg {
		m := envs.NewMinio()
		copy(item, m)
		items[key] = m
	}
	return items
}

func NewNeo4j(cfg map[string]*Neo4jConfig) map[string]*envs.Neo4j {
	items := map[string]*envs.Neo4j{}
	for key, item := range cfg {
		neo4j := envs.NewNeo4j()
		copy(item, neo4j)
		items[key] = neo4j
	}
	return items
}

func NewMySQL(cfg map[string]*MySqlConfig) map[string]*envs.MySql {
	items := map[string]*envs.MySql{}
	for key, item := range cfg {
		mysql := envs.NewMySQL()
		copy(item, mysql)
		items[key] = mysql
	}
	return items
}

func NewMongo(cfg map[string]*MongoConfig) map[string]*envs.Mongo {
	items := map[string]*envs.Mongo{}
	for key, item := range cfg {
		mysql := &envs.Mongo{}
		copy(item, mysql)
		items[key] = mysql
	}
	return items
}

func NewDapr(cfg *DaprConfig) *envs.Dapr {
	dapr := envs.NewDapr()
	copy(cfg, dapr)
	return dapr
}

/*
func NewAuth(cfg *AuthConfig) *env.Auth {
	auth := env.NewAuth()
	copy(cfg, auth)
	return auth
}
*/

func NewApp(cfg *AppConfig) *envs.App {
	app := envs.NewApp()
	copy(cfg, app)
	return app
}

func NewLog(cfg *LogConfig) *envs.Log {
	log := envs.NewLog()
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
