package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	log "github.com/sirupsen/logrus"
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
	Name  string                  `yaml:"-"`
	App   AppConfig               `yaml:"app"`
	Log   LogConfig               `yaml:"log"`
	Dapr  DaprConfig              `yaml:"dapr"`
	Mongo map[string]*MongoConfig `yaml:"mongo"`
	Neo4j map[string]*Neo4jConfig `yaml:"neo4j"`
	Mysql map[string]*MySqlConfig `yaml:"mysql"`
	Minio map[string]*MinioConfig `yaml:"minio"`
	Redis map[string]*RedisConfig `yaml:"redis"`
}

type AppConfig struct {
	AppId    string            `yaml:"id"`
	HttpHost string            `yaml:"httpHost"`
	HttpPort int               `yaml:"httpPort"`
	RootUrl  string            `yaml:"rootUrl"`
	Values   map[string]string `yaml:"values"`
}

type DaprConfig struct {
	Host               *string                `yaml:"host"`
	HttpPort           *int64                 `yaml:"httpPort"`
	GrpcPort           *int64                 `yaml:"grpcPort"`
	MaxCallRecvMsgSize *int64                 `yaml:"maxCallRecvMsgSize"` //dapr数据包大小，单位M
	EventStores        map[string]*EventStore `yaml:"eventStores"`
	Actor              ActorConfig            `yaml:"actor"`
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
	Level string `yaml:"level"`
	File  string `yaml:"file"`
	level applog.Level
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

	if e.Log.Level != "" {
		l, err := logs.ParseLevel(e.Log.Level)
		if err != nil {
			return err
		}
		e.Log.level = l
	}

	initLogs(e.Log.level)

	if e.Dapr.Host == nil {
		var value = "localhost"
		e.Dapr.Host = e.GetEnvString("DAPR_HOST", &value)
	}

	if e.Dapr.HttpPort == nil {
		var value int64 = 3500
		e.Dapr.HttpPort = e.GetEnvInt("DAPR_HTTP_PORT", &value)
	}

	if e.Dapr.GrpcPort == nil {
		var value int64 = 50001
		e.Dapr.GrpcPort = e.GetEnvInt("DAPR_GRPC_PORT", &value)
	}

	if len(e.Dapr.EventStores) > 0 {
		for compName, es := range e.Dapr.EventStores {
			if es.CompName == "" {
				es.CompName = compName
			}
			if len(es.PubsubName) == 0 {
				return errors.ErrorOf("config env:%s  Dapr.EventStores.%s pubsub is null", name, compName)
			}
		}
	}

	e.Dapr.Actor.init()

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

func (d DaprConfig) GetHost() string {
	if d.Host == nil {
		return ""
	}
	return *d.Host
}

func (d DaprConfig) GetHttpPort() int64 {
	if d.HttpPort == nil {
		return 0
	}
	return *d.HttpPort
}

func (d DaprConfig) GetGrpcPort() int64 {
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

func (c *Config) GetEnvConfig(env string) (*EnvConfig, error) {
	envConfig, ok := c.Envs[env]
	if !ok {
		return nil, errors.New("not found env: " + env)
	}

	if envConfig != nil {
		log.Infoln(fmt.Sprintf("ctype=app; appId=%s; env=%s;", envConfig.App.AppId, env))
		log.Infoln(fmt.Sprintf("ctype=app; httpHost=%s; httpPort=%d; httpRootUrl=%s;", envConfig.App.HttpHost, envConfig.App.HttpPort, envConfig.App.RootUrl))
		log.Infoln(fmt.Sprintf("ctype=dapr; daprHost=%s; daprHttpPort=%d; daprGrpcPort=%d;", envConfig.Dapr.GetHost(), envConfig.Dapr.GetHttpPort(), envConfig.Dapr.GetGrpcPort()))
		log.Infoln(fmt.Sprintf("ctype=eventStores; length=%v;", len(envConfig.Dapr.EventStores)))
		return envConfig, nil
	}

	return nil, NewEnvTypeError(fmt.Sprintf("error config env is \"%s\". choose one of: [dev, test, prod]", env))
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
