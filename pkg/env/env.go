package env

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	logs2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
)

type Config struct {
	Env  string          `yaml:"env"`
	Envs map[string]*Env `yaml:"envs"`
}

type Env struct {
	Name      string               `yaml:"-" json:"name"`
	App       *App                 `yaml:"app" json:"app"`
	Log       *Log                 `yaml:"log" json:"log"`
	Dapr      *Dapr                `yaml:"dapr" json:"dapr"`
	Resources map[string]*Resource `yaml:"resources" json:"resources"`
	Mongo     map[string]*Mongo    `yaml:"mongo" json:"mongo"`
	Neo4j     map[string]*Neo4j    `yaml:"neo4j" json:"neo4J"`
	Mysql     map[string]*MySql    `yaml:"mysql" json:"mysql"`
	Minio     map[string]*Minio    `yaml:"minio" json:"minio"`
	Redis     map[string]*Redis    `yaml:"redis" json:"redis"`
	Fs        []map[string]any     `yaml:"fs" json:"fs"`
	Auth      *Auth                `yaml:"auth" json:"auth"`
	Fsm       *fsm.Manager         `yaml:"-" json:"-"`
	dbs       map[string]DBItem
}

var _env *Env

func SetEnv(env *Env) {
	_env = env
	if env != nil {
		SetEnvName(env.Name)

		//envConfig.Dapr.Server.Config = AbsFileName(envConfig.Dapr.Server.Config)
		//envConfig.Dapr.Server.LogFile = AbsFileName(envConfig.Dapr.Server.LogFile)
		//envConfig.Dapr.Server.ComponentsPath = AbsFileName(envConfig.Dapr.Server.ComponentsPath)
	}
}

func GetEnv() *Env {
	return _env
}

func NewEnv() *Env {
	return &Env{
		Name:      "",
		App:       NewApp(),
		Log:       NewLog(),
		Dapr:      NewDapr(),
		Resources: map[string]*Resource{},
		Mongo:     map[string]*Mongo{},
		Neo4j:     map[string]*Neo4j{},
		Mysql:     map[string]*MySql{},
		Minio:     map[string]*Minio{},
		Redis:     map[string]*Redis{},
		Fs:        []map[string]any{},
		Auth:      &Auth{},
		Fsm:       fsm.NewManager(),
		dbs:       map[string]DBItem{},
	}
}

func (env *Env) Init() {
	if env == nil {
		return
	}
	for _, m := range env.Fs {
		for k, v := range m {
			if s, ok := v.(string); ok {
				m[k] = ReplaceSysValues(s)
			}
		}
	}
	initLog(env)
	initApp(env)

	InitDBMongo(env)
	InitDBMySql(env)
	InitDBNeo4j(env)

	initMinio(env)
	initDapr(env)
	initResources(env)

}

func (env *Env) NewFsManager(fs []map[string]any) *fsm.Manager {
	if len(env.Fs) != 0 {
		fsManager, err := fsm.NewManagerWithConfigs(fs, env.App.HServer.SrcName)
		if err != nil {
			panic(errors.New("fs.NewManagerWithConfigs() err: %s", err.Error()))
		}
		return fsManager
	}
	return fsm.NewManager()
}

func (env *Env) GetProdMode() bool {
	return env.App.ProdMode
}

func (env *Env) AddDB(dbItem DBItem) {
	dbKey := dbItem.GetDBKey()
	_, ok := env.dbs[dbKey]
	if ok {
		panic(errors.New("db \"%s\" already exists", dbKey))
	}
	env.dbs[dbKey] = dbItem
}

func (env *Env) AddNeo4j(dbCfg *Neo4j) {
	if dbCfg == nil {
		panic(errors.New("dbCfg is nil"))
	}
	env.Neo4j[dbCfg.DbKey] = dbCfg
}

func (env *Env) AddMongo(dbCfg *Mongo) {
	if dbCfg == nil {
		panic(errors.New("dbCfg is nil"))
	}
	env.Mongo[dbCfg.DbKey] = dbCfg
}

func (env *Env) AddMySql(dbCfg *MySql) {
	if dbCfg == nil {
		panic(errors.New("dbCfg is nil"))
	}
	env.Mysql[dbCfg.DbKey] = dbCfg
}

func (env *Env) CloseDB(ctx context.Context) error {
	for _, d := range env.dbs {
		_ = d.CloseDB(ctx)
	}
	return nil
}

func (env *Env) GetDB(dbKey string) DBItem {
	item := env.dbs[dbKey]
	return item
}

func (env *Env) GetField(key string) any {
	return reflectutils.GetField(env, key)
}

func GetDB(dbKey string) DBItem {
	return _env.GetDB(dbKey)
}

func CloseDB(ctx context.Context) error {
	return _env.CloseDB(ctx)
}

func GetDBDefault() DBItem {
	env := GetEnv()
	for _, item := range env.dbs {
		return item
	}
	return nil
}

func GetMongoByKey(dbKey string) (*mongodb.MongoDB, bool) {
	item := GetDB(dbKey)
	if item == nil {
		return nil, false
	}
	mdb := item.GetMongo()
	return mdb, mdb != nil
}

func GetAppValue(name string) (any, error) {
	var err error
	v, ok := _env.App.Meta[name]
	if !ok {
		err = errors.New(fmt.Sprintf("配置变量%s不存在", name))
	}
	return v, err
}

func GetAppValues() map[string]any {
	return _env.App.Meta
}

func GetDaprHost() string {
	return _env.Dapr.GetHost()
}

func GetDaprHttpPort() int64 {
	return _env.Dapr.GetHttpPort()
}

func GetDaprGrpcPort() int64 {
	return _env.Dapr.GetGrpcPort()
}

func GetAppId() string {
	return _env.App.AppId
}

func GetAppName() string {
	return _env.App.AppName
}

func GetAppHttpHost() string {
	return _env.App.HttpHost
}

func GetHttpInvoke(appId string) string {
	return fmt.Sprintf("http://%s:%v/v1.0/invoke/%v/method/", GetDaprHost(), GetDaprHttpPort(), appId)
}

func GetHttpsInvoke(appId string) string {
	return fmt.Sprintf("https://%s:%v/v1.0/invoke/%v/method/", GetDaprHost(), GetDaprHttpPort(), appId)
}

func GetLogger() logs2.Logger {
	return logs2.GetLogger()
}
