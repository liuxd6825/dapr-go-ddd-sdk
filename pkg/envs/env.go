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
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
	"github.com/minio/minio-go/v7"
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
	Services  map[string]Service   `yaml:"services" json:"services"`
	DBs       map[string]Metadata  `yaml:"dbs" json:"dbs"`
	OSS       map[string]Metadata  `yaml:"oss" json:"oss"`
	Graphdb   map[string]Metadata  `yaml:"graphdb" json:"graphdb"`
	IMDB      map[string]Metadata  `yaml:"imdb" json:"imdb"`
	Fs        []map[string]any     `yaml:"fs" json:"fs"`
	Fsm       *fsm.Manager         `yaml:"-" json:"-"`
	oss       *OssEnv              `yaml:"-" json:"-"` // 对象存储服务
	dbs       *DbsEnv              `yaml:"-" json:"-"` // 系统数据库或json数据库
	graphdb   *GraphEnv            `yaml:"-" json:"-"` // 图数据库
	imdb      *IMDBEnv             `yaml:"-" json:"-"` // 内存数据库
}

type Metadata map[string]any

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
		DBs:       map[string]Metadata{},
		OSS:       map[string]Metadata{},
		Graphdb:   map[string]Metadata{},
		IMDB:      map[string]Metadata{},
		Fs:        []map[string]any{},
		Fsm:       fsm.NewManager(),
		dbs:       newDbsEnv(),
		oss:       newOssEnv(),
		graphdb:   newGraphEnv(),
		imdb:      newIMDBEnv(),
	}
}

func (env *Env) Init() *Env {
	if env == nil {
		panicError(errors.New("env is nil"))
		return nil
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

	initOSS(env)
	initDB(env)
	initGraphDB(env)
	initIMDB(env)

	env.Fsm = env.NewFsManager(env.Fs)
	return env
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

func (env *Env) AddOSS(fsItem OSSItem) {
	fsKey := fsItem.GetOssKey()
	_, ok := env.oss.items[fsKey]
	if ok {
		panicError(errors.New("fileSystem \"%s\" already exists", fsKey))
	}
}

func (env *Env) AddDB(dbItem DBItem) {
	dbKey := dbItem.GetDBKey()
	_, ok := env.dbs.items[dbKey]
	if ok {
		panicError(errors.New("db \"%s\" already exists", dbKey))
	}
	env.dbs.items[dbKey] = dbItem
}

func (env *Env) AddIMDB(dbItem IMDBItem) {
	dbKey := dbItem.GetDBKey()
	_, ok := env.imdb.items[dbKey]
	if ok {
		panicError(errors.New("db \"%s\" already exists", dbKey))
	}
	env.imdb.items[dbKey] = dbItem
}

func (env *Env) CloseDB(ctx context.Context) error {
	for _, d := range env.dbs.items {
		_ = d.CloseDB(ctx)
	}
	return nil
}

func (env *Env) GetDB(dbKey string) DBItem {
	dbKey = env.GetDBKeyValue(dbKey)
	item := env.dbs.items[dbKey]
	return item
}

func (env *Env) GetFs(dbKey string) OSSItem {
	dbKey = strings.ToLower(dbKey)
	item := env.oss.items[dbKey]
	return item
}

func (env *Env) GetMinioConfig(fsKey string) *Minio {
	m, ok := env.oss.minio[fsKey]
	if !ok {
		panicError(errors.New("minio fsKey %s not found", fsKey))
	}
	return m
}

func (env *Env) GetMinioClient(fsKey string) *minio.Client {
	item, ok := env.oss.minio[fsKey]
	if !ok {
		panicError(errors.New("minio fsKey %s not found", fsKey))
	}
	return item.Client
}

func (env *Env) GetHugeConfig(graphdbKey string) *Huge {
	cfg, ok := env.graphdb.Huge[graphdbKey]
	if !ok {
		panicError(errors.New("huge graphdbKey %s not found", graphdbKey))
	}
	return cfg
}

func (env *Env) GetRedisConfig(dbKey string) *Redis {
	cfg, ok := env.imdb.redis[dbKey]
	if !ok {
		panicError(errors.New("redis imdbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetNeo4jConfig(dbKey string) *Neo4j {
	cfg, ok := env.graphdb.Neo4j[dbKey]
	if !ok {
		panicError(errors.New("neo4j graphdbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetElasticConfig(dbKey string) *Elastic {
	cfg, ok := env.dbs.elastic[dbKey]
	if !ok {
		panicError(errors.New("elastic dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetMongoConfig(dbKey string) *Mongo {
	cfg, ok := env.dbs.mongo[dbKey]
	if !ok {
		panicError(errors.New("mongo dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetPostgresConfig(dbKey string) *Postgres {
	cfg, ok := env.dbs.postgres[dbKey]
	if !ok {
		panicError(errors.New("postgres dbKey %s not found", dbKey))
	}
	return cfg
}

func (env *Env) GetNeo4jClient(dbKey string) neo4j.DriverWithContext {
	dbKey = env.GetDBKeyValue(dbKey)
	item, isHas := env.dbs.items[dbKey]
	if !isHas {
		panicError(errors.New("neo4j graphdbKey \"%s\" not found", dbKey))
	}
	if item.GetDBType() != DBTypeNeo4j {
		panicError(errors.New("graphdbKey %s is not neo4j, ", dbKey))
	}
	return item.GetNeo4jClient()
}

func (env *Env) GetHugeClient(dbKey string) *hugegraph.CommonClient {
	dbKey = env.GetDBKeyValue(dbKey)
	item := env.dbs.items[dbKey]
	if item == nil {
		panicError(errors.New("huge graphdbKey \"%s\" not found", dbKey))
	}
	return item.GetHugeClient()
}

func (env *Env) GetRedisClient(key string) *redis.Client {
	item := env.imdb.items[key]
	if item == nil {
		panicError(errors.New("redis imdbKey \"%s\" not found", key))
	}
	if item.GetIMDBType() != IMDBTypeRedis {
		panicError(errors.New("imdbType is invalid, ", item.GetIMDBType()))
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

func GetMinioByKey(fsKey string) (*Minio, bool) {
	item := GetEnv()
	if item == nil {
		return nil, false
	}
	mdb := item.GetMinioConfig(fsKey)
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
