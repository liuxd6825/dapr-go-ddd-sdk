package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Env  string                `yaml:"env"`
	Envs map[string]*EnvConfig `yaml:"envs"`
}

type EnvConfig struct {
	Name string     `yaml:"-"`
	App  AppConfig  `yaml:"app"`
	Log  LogConfig  `yaml:"log"`
	Dapr DaprConfig `yaml:"dapr"`

	Resources map[string]*ResourceConfig `yaml:"resources"`
	Mongo     map[string]*MongoConfig    `yaml:"mongo"`
	Neo4j     map[string]*Neo4jConfig    `yaml:"neo4j"`
	Mysql     map[string]*MySqlConfig    `yaml:"mysql"`
	Minio     map[string]*MinioConfig    `yaml:"minio"`
	Redis     map[string]*RedisConfig    `yaml:"redis"`
}

type AppConfig struct {
	AppId     string            `yaml:"id"`
	AppName   string            `yaml:"name"`
	HttpHost  string            `yaml:"httpHost"`
	HttpPort  int               `yaml:"httpPort"`
	RootUrl   string            `yaml:"rootUrl"`
	CPU       *int              `yaml:"cpu"`
	Memory    *string           `yaml:"memory"`
	Values    map[string]string `yaml:"values"`
	AuthToken string            `yaml:"authToken"`
}

type ResourceConfig struct {
	Namespace string            `yaml:"namespace"`
	Name      string            `yaml:"name"`
	Type      string            `yaml:"type"`
	URI       string            `yaml:"uri"`
	Metadata  map[string]string `yaml:"metadata"`
}

type DaprConfig struct {
	Host                *string                `yaml:"host"`
	HttpPort            *int64                 `yaml:"httpPort"`
	GrpcPort            *int64                 `yaml:"grpcPort"`
	MaxCallRecvMsgSize  *int                   `yaml:"maxCallRecvMsgSize"` //dapr数据包大小，单位M
	MaxIdleConns        *int                   `yaml:"maxIdleConns"`
	MaxIdleConnsPerHost *int                   `yaml:"maxIdleConnsPerHost"`
	IdleConnTimeout     *int                   `yaml:"idleConnTimeout"`
	EventStores         map[string]*EventStore `yaml:"eventStores"`
	Actor               ActorConfig            `yaml:"actor"`
	Server              DaprServerConfig       `yaml:"server"`
}

// DaprServerConfig dapr服务端参数
type DaprServerConfig struct {
	Start                bool   `yaml:"start"` //是否启动Daprd
	LogLevel             string `yaml:"logLevel"`
	EnableMetrics        bool   `yaml:"enableMetrics"`
	Config               string `yaml:"config"`
	ComponentsPath       string `yaml:"componentsPath"`
	PlacementHostAddress string `yaml:"placementHostAddress"`
}

type ActorConfig struct {
	ActorIdleTimeout       string `yaml:"actorIdleTimeout"`
	ActorScanInterval      string `yaml:"actorScanInterval"`
	DrainOngingCallTimeout string `yaml:"drainOngoingCallTimeout"`
	DrainBalancedActors    bool   `yaml:"drainRebalancedActors"`
}

type EventStore struct {
	CompName   string `yaml:"name"`   // Dapr EventStarge 组件名称
	PubsubName string `yaml:"pubsub"` // Dapr Pubsub 组件名称
}

type LogConfig struct {
	Level     string `yaml:"level"`
	SaveDays  int    `yaml:"saveDays"`  //日志保存的天数
	SplitHour int    `yaml:"splitHour"` //文件分隔时间，单位小时
	level     logs.Level
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigByFile(fileName string) (*Config, error) {
	//rootPath, _ := os.Getwd()
	//_ = fmt.Sprintf("%s/%s", rootPath, fileName)
	filename := fileName
	if strings.HasPrefix(filename, "${search}") {
		slist := strings.Split(filename, "/")
		slist = slist[1:]
		v, ok, err := searchConfigFile(".", slist[0], slist[1])
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.New(fileName)
		}
		filename = v
	}
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	for name, env := range config.Envs {
		if err := env.Init(name); err != nil {
			return nil, err
		}
	}
	return &config, nil
}

func (e *EnvConfig) Init(name string) error {
	e.Name = name
	if len(e.App.HttpHost) == 0 {
		e.App.HttpHost = "0.0.0.0"
	}

	// init log
	level := logs.ErrorLevel
	if e.Log.Level != "" {
		l, err := logs.ParseLevel(e.Log.Level)
		if err != nil {
			return err
		}
		level = l
	}
	if e.Log.SaveDays <= 0 {
		e.Log.SaveDays = 30
	}
	if e.Log.SplitHour <= 0 {
		e.Log.SplitHour = 24
	}
	e.Log.level = level

	//初始化日志
	if err := initLogs(level, e.Log.SaveDays, e.Log.SplitHour); err != nil {
		return err
	}

	//初始化Dapr
	if err := e.Dapr.init(e); err != nil {
		return err
	}

	return nil
}

func (e *EnvConfig) GetEnvInt(envName string, defValue *int64) *int64 {
	value, ok := os.LookupEnv(envName)
	if !ok {
		return defValue
	}
	parseInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return &parseInt
}

func (e *EnvConfig) GetEnvString(envName string, defValue *string) *string {
	value, ok := os.LookupEnv(envName)
	if !ok {
		return defValue
	}
	return &value
}

func (d *DaprConfig) GetHost() string {
	if d.Host == nil {
		return ""
	}
	return *d.Host
}

func (d *DaprConfig) GetHttpPort() int64 {
	if d.HttpPort == nil {
		return 0
	}
	return *d.HttpPort
}

func (d *DaprConfig) GetGrpcPort() int64 {
	if d.GrpcPort == nil {
		return 0
	}
	return *d.GrpcPort
}

func (l *LogConfig) GetLevel() applog.Level {
	return l.level
}

func (c *ActorConfig) init() {
	if c.ActorIdleTimeout == "" {
		c.ActorIdleTimeout = "1h"
	}
	if c.ActorScanInterval == "" {
		c.ActorScanInterval = "30s"
	}
	if c.DrainOngingCallTimeout == "" {
		c.DrainOngingCallTimeout = "5m"
	}
}

func (c *DaprConfig) init(e *EnvConfig) error {
	if c.Host == nil {
		var value = "localhost"
		c.Host = e.GetEnvString("DAPR_HOST", &value)
	}

	if e.Dapr.HttpPort == nil {
		var value int64 = 3500
		c.HttpPort = e.GetEnvInt("DAPR_HTTP_PORT", &value)
	}

	if e.Dapr.GrpcPort == nil {
		var value int64 = 50001
		c.GrpcPort = e.GetEnvInt("DAPR_GRPC_PORT", &value)
	}

	if c.MaxCallRecvMsgSize == nil {
		val := dapr.GetMaxCallRecvMsgSize()
		c.MaxCallRecvMsgSize = &val
	}

	if c.MaxIdleConnsPerHost == nil {
		val := dapr.DefaultMaxIdleConnsPerHost
		c.MaxIdleConns = &val
	}

	if c.IdleConnTimeout == nil {
		val := dapr.DefaultIdleConnTimeout
		c.IdleConnTimeout = &val
	}

	if c.MaxIdleConns == nil {
		val := dapr.DefaultMaxIdleConns
		c.MaxIdleConnsPerHost = &val
	}

	if len(c.EventStores) > 0 {
		for compName, es := range e.Dapr.EventStores {
			if es.CompName == "" {
				es.CompName = compName
			}
			if len(es.PubsubName) == 0 {
				return errors.ErrorOf("config env:%s  Dapr.EventStores.%s pubsub is null", e.Name, compName)
			}
		}
	}

	e.Dapr.Actor.init()

	return nil

}

func (c *Config) GetEnvConfig(env string) (*EnvConfig, error) {
	envConfig, ok := c.Envs[env]
	if !ok {
		return nil, errors.New("not found env: " + env)
	}

	if envConfig != nil {
		return envConfig, nil
	}

	return nil, NewEnvTypeError(fmt.Sprintf("error config env is \"%s\". choose one of: [dev, test, prod]", env))
}

func initResources(resCfg map[string]*ResourceConfig) error {
	if resCfg == nil {
		return nil
	}

	for k, v := range resCfg {
		v.Name = k
	}
	return nil
}

func searchConfigFile(path, configName string, fileName string) (string, bool, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", false, err
	}
	if len(files) <= 0 {
		return "", false, nil
	}
	for _, file := range files {
		name := file.Name()
		if file.IsDir() && name == configName {
			list, err := ioutil.ReadDir(path + "/" + file.Name())
			if err != nil {
				return "", false, err
			}
			for _, item := range list {
				if item.Name() == fileName {
					return fmt.Sprintf("%v/%v/%v", path, file.Name(), item.Name()), true, nil
				}
			}
		}
	}

	return searchConfigFile(path+"/..", configName, fileName)
}
