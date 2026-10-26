package envs

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	logs2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/boolutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
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
	/*
		Mongo     map[string]*Mongo    `yaml:"mongo" json:"mongo"`
		Neo4j     map[string]*Neo4j    `yaml:"neo4jdao" json:"neo4J"`
		Mysql     map[string]*MySql    `yaml:"mysql" json:"mysql"`
		Nebula    map[string]*Nebula   `yaml:"nebula" json:"nebula"`
		Elastic   map[string]*Elastic  `yaml:"elastic" json:"elastic"`

		Redis     map[string]*Redis    `yaml:"redis" json:"redis"`
	*/

	Minio    map[string]*Minio  `yaml:"minio" json:"minio"`
	DBs      map[string]DBInfo  `yaml:"dbs" json:"dbs"`
	Services map[string]Service `yaml:"services" json:"services"`
	//Auth      *Auth                `yaml:"auth" json:"auth"`
	Fs        []map[string]any `yaml:"fs" json:"fs"`
	Fsm       *fsm.Manager     `yaml:"-" json:"-"`
	dbs       map[string]DBItem
	dbsConfig *DbsConfig
}

type DbsConfig struct {
	Mongo    map[string]*Mongo
	Neo4j    map[string]*Neo4j
	Mysql    map[string]*MySql
	Doris    map[string]*Doris
	Postgres map[string]*Postgres
	Nebula   map[string]*Nebula
	Elastic  map[string]*Elastic
	Huge     map[string]*Huge
	Redis    map[string]*Redis
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

type DBInfo map[string]any

func NewEnv() *Env {
	return &Env{
		Name:      "",
		App:       NewApp(),
		Log:       NewLog(),
		Dapr:      NewDapr(),
		Resources: map[string]*Resource{},
		DBs:       map[string]DBInfo{},
		Minio:     map[string]*Minio{},
		Fs:        []map[string]any{},
		Fsm:       fsm.NewManager(),
		dbs:       map[string]DBItem{},
		dbsConfig: newDbsConfig(),
	}
}

func newDbsConfig() *DbsConfig {
	return &DbsConfig{
		Mongo:    map[string]*Mongo{},
		Neo4j:    map[string]*Neo4j{},
		Mysql:    map[string]*MySql{},
		Doris:    map[string]*Doris{},
		Postgres: map[string]*Postgres{},
		Nebula:   map[string]*Nebula{},
		Elastic:  map[string]*Elastic{},
		Huge:     map[string]*Huge{},
		Redis:    map[string]*Redis{},
	}
}

func (env *Env) Init() *Env {
	if env == nil {
		panicError(errors.New("env is nil"))
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
	initResources(env)
	initDapr(env)
	initMinio(env)
	InitDB(env)

	env.Fsm = env.NewFsManager(env.Fs)
	return env
}

func InitDB(env *Env) {
	errList := errors.NewErrors()
	for dbKey, info := range env.DBs {
		dbType, isHas := info["dbType"]
		info["dbKey"] = strings.ToLower(dbKey)
		if !isHas {
			errList.AddError(errors.New("%s dbType is required", dbKey))
		}
		enable := boolutils.GetBool(info["enable"], true)
		if !enable {
			continue
		}

		switch dbType {
		case DBTypeDoris.String():
			db, err := NewDorisWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBDoris(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Doris[db.DbKey] = db
			break
		case DBTypePostgres.String():
			db, err := NewPostgresWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBPostgres(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Postgres[db.DbKey] = db
			break
		case DBTypeMySQL.String():
			db, err := NewMySQLWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBMySql(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Mysql[db.DbKey] = db
			break
		case DBTypeMongoDB.String():
			db, err := NewMongoWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBMongo(env, db); err != nil {
				errList.AddError(err)
			}

			env.dbsConfig.Mongo[db.DbKey] = db
			break
		case DBTypeNeo4j.String():
			db, err := NewNeo4jWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBNeo4j(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Neo4j[db.DbKey] = db
			break
		case DBTypeElastic.String():
			db, err := NewElasticWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBElastic(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Elastic[db.DbKey] = db
			break
		case DBTypeRedis.String():
			db, err := NewRedisWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBRedis(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Redis[db.DbKey] = db
			break

		case DBTypeHuge.String():
			db, err := NewHugeWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBHuge(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Huge[db.DbKey] = db
			break
		case DBTypeNebula.String():
			db, err := NewNebulaWithMap(info)
			if err != nil {
				errList.AddError(errors.New("dbKey %s error: %s", dbKey, err.Error()))
				continue
			}
			if err := AddDBNebula(env, db); err != nil {
				errList.AddError(err)
			}
			env.dbsConfig.Nebula[db.DbKey] = db
			break
		default:
			errList.AddError(errors.New("dbType is invalid, ", dbKey))
		}

	}

	if errList.HasError() {
		panicError(errList)
	}
}

func (env *Env) NewFsManager(fs []map[string]any) *fsm.Manager {
	if len(env.Fs) != 0 {
		fsManager, err := fsm.NewManagerWithConfigs(fs, env.App.HServer.SrcName)
		if err != nil {
			panicError(errors.New("fs.NewManagerWithConfigs() err: %s", err.Error()))
		}
		return fsManager
	}
	return fsm.NewManager()
}

func (env *Env) GetServiceConfig(serviceKey string, cfg any) {
	meta := env.Services[serviceKey]
	err := copyObject(meta, cfg)
	if err != nil {
		panicError(errors.New("%s config %s error: %s ", serviceKey, err.Error()))
	}
}

func (env *Env) GetProdMode() bool {
	return env.App.ProdMode
}

func panicError(err error) {
	file := reflectutils.CallerInfo(4)
	panic(err.Error() + ": " + file)
}

func (env *Env) AddDB(dbItem DBItem) {
	dbKey := dbItem.GetDBKey()
	_, ok := env.dbs[dbKey]
	if ok {
		panicError(errors.New("db \"%s\" already exists", dbKey))
	}
	env.dbs[dbKey] = dbItem
}

/*
func (env *Env) AddNeo4j(neo4jCfg *Neo4j) {
	if neo4jCfg == nil {
		panicError(errors.New("dbCfg is nil"))
	}
	env.Neo4j[neo4jCfg.DbKey] = neo4jCfg
}

func (env *Env) AddMongo(mongoCfg *Mongo) {
	if mongoCfg == nil {
		panicError(errors.New("mongoCfg is nil"))
	}
	env.Mongo[mongoCfg.DbKey] = mongoCfg
}

func (env *Env) AddMySql(mysqlCfg *MySql) {
	if mysqlCfg == nil {
		panicError(errors.New("mysqlCfg is nil"))
	}
	env.Mysql[mysqlCfg.DbKey] = mysqlCfg
}

func (env *Env) AddHuge(hugeCfg *Huge) {
	if hugeCfg == nil {
		panicError(errors.New("hugeCfg is nil"))
	}
	env.Huge[hugeCfg.DbKey] = hugeCfg
}

func (env *Env) AddMinio(minoCfg *Minio) {
	if minoCfg == nil {
		panicError(errors.New("minoCfg is nil"))
	}
	env.Minio[minoCfg.Name] = minoCfg
}*/

func (env *Env) CloseDB(ctx context.Context) error {
	for _, d := range env.dbs {
		_ = d.CloseDB(ctx)
	}
	return nil
}

func (env *Env) GetDB(dbKey string) DBItem {
	dbKey = env.GetDBKeyValue(dbKey)
	item := env.dbs[dbKey]
	return item
}

func (env *Env) GetHugeConfig(dbKey string) *Huge {
	cfg, ok := env.dbsConfig.Huge[dbKey]
	if !ok {
		panicError(errors.New("huge dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetRedisConfig(dbKey string) *Redis {
	cfg, ok := env.dbsConfig.Redis[dbKey]
	if !ok {
		panicError(errors.New("redis dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetNeo4jConfig(dbKey string) *Neo4j {
	cfg, ok := env.dbsConfig.Neo4j[dbKey]
	if !ok {
		panicError(errors.New("neo4j dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetElasticConfig(dbKey string) *Elastic {
	cfg, ok := env.dbsConfig.Elastic[dbKey]
	if !ok {
		panicError(errors.New("elastic dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetMongoConfig(dbKey string) *Mongo {
	cfg, ok := env.dbsConfig.Mongo[dbKey]
	if !ok {
		panicError(errors.New("mongo dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetPostgresConfig(dbKey string) *Postgres {
	cfg, ok := env.dbsConfig.Postgres[dbKey]
	if !ok {
		panicError(errors.New("postgres dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetNeo4jClient(dbKey string) neo4j.DriverWithContext {
	dbKey = env.GetDBKeyValue(dbKey)
	item, isHas := env.dbs[dbKey]
	if !isHas {
		panicError(errors.New("neo4j dbKey \"%s\" not found", dbKey))
	}
	if item.GetDBType() != DBTypeNeo4j {
		panicError(errors.New("dbKey %s is not neo4j, ", dbKey))
	}
	return item.GetNeo4jClient()
}

func (env *Env) GetHugeClient(dbKey string) *hugegraph.CommonClient {
	dbKey = env.GetDBKeyValue(dbKey)
	item := env.dbs[dbKey]
	if item == nil {
		panicError(errors.New("huge dbKey \"%s\" not found", dbKey))
	}
	return item.GetHugeClient()
}

func (env *Env) GetRedisClient(key string) *redis.Client {
	item := env.GetDB(key)
	if item == nil {
		panicError(errors.New("redis dbKey \"%s\" not found", key))
	}
	if item.GetDBType() != DBTypeRedis {
		panicError(errors.New("dbType is invalid, ", item.GetDBType()))
	}
	return item.GetRedisClient()
}

func (env *Env) GetField(key string) any {
	return reflectutils.GetField(env, key)
}

func (env *Env) GetDBKeyValue(dbKey string) string {
	if strings.HasPrefix(dbKey, "$") {
		dbKey = dbKey[1:len(dbKey)]
		gp.Try(func() error {
			if keyVal, ok := env.App.Meta[dbKey]; ok {
				dbKey = keyVal.(string)
			} else {
				panicError(errors.New("invalid db key", dbKey))
			}
			return nil
		}).Catch(func(e error) {
			panicError(errors.New("dbKey must start with $%s", dbKey))
		})
	}
	return dbKey
}

func (env *Env) GetMinioByKey(dbKey string) (*Minio, bool) {
	mdb := env.Minio[dbKey]
	return mdb, mdb != nil
}

func (env *Env) GetMongoByKey(dbKey string) *mongo.Client {
	item := env.GetDB(dbKey)
	if item == nil {
		panicError(errors.New("db \"%s\" not found", dbKey))
	}
	if item.GetDBType() != DBTypeMongoDB {
		panicError(errors.New("dbType is invalid, ", item.GetDBType()))
	}
	return item.GetMongoClient()
}

//

func GetDB(dbKey string) DBItem {
	dbKey = _env.GetDBKeyValue(dbKey)
	return _env.GetDB(dbKey)
}

func CloseDB(ctx context.Context) error {
	return _env.CloseDB(ctx)
}

/*func GetDBDefault() DBItem {
	env := GetEnv()
	for _, item := range env.dbs {
		return item
	}
	return nil
}
*/

func GetMongoClient(dbKey string) *mongo.Client {
	item := GetDB(dbKey)
	if item == nil {
		panicError(errors.New("db \"%s\" not found", dbKey))
	}
	if item.GetDBType() != DBTypeMongoDB {
		panicError(errors.New("dbType is invalid, ", item.GetDBType()))
	}
	return item.GetMongoClient()
}

func GetMongo(dbKey string) *mongodb.MongoDB {
	item := GetDB(dbKey)
	if item == nil {
		panicError(errors.New("db \"%s\" not found", dbKey))
	}
	if item.GetDBType() != DBTypeMongoDB {
		panicError(errors.New("dbType is invalid, ", item.GetDBType()))
	}
	return item.GetMongoDB()
}

func GetMinioByKey(dbKey string) (*Minio, bool) {
	item := GetEnv()
	if item == nil {
		return nil, false
	}
	mdb := item.Minio[dbKey]
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
