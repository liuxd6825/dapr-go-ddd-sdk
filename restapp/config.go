package restapp

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	EnvType string                `yaml:"envType"`
	Test    *EnvConfig            `yaml:"test"`
	Dev     *EnvConfig            `yaml:"dev"`
	Prod    *EnvConfig            `yaml:"prod"`
	Envs    map[string]*EnvConfig `yaml:"envs"`
}

type EnvConfig struct {
	App   AppConfig   `yaml:"app"`
	Log   LogConfig   `yaml:"log"`
	Dapr  DaprConfig  `yaml:"dapr"`
	Mongo MongoConfig `yaml:"mongo"`
}

func (e *EnvConfig) Init() error {

	if len(e.App.HttpHost) == 0 {
		e.App.HttpHost = "0.0.0.0"
	}

	if e.Log.Level != "" {
		l, err := applog.NewLevel(e.Log.Level)
		if err != nil {
			return err
		}
		e.Log.level = l
	}

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

type AppConfig struct {
	AppId    string `yaml:"id"`
	HttpHost string `yaml:"httpHost"`
	HttpPort int    `yaml:"httpPort"`
	RootUrl  string `yaml:"rootUrl"`
}

type DaprConfig struct {
	Host     *string  `yaml:"host"`
	HttpPort *int64   `yaml:"httpPort"`
	GrpcPort *int64   `yaml:"grpcPort"`
	Pubsubs  []string `yaml:"pubsubs,flow"`
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

type LogConfig struct {
	Level string `yaml:"level"`
	level applog.Level
}

func (l *LogConfig) GetLevel() applog.Level {
	return l.level
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigByFile(fileName string) (*Config, error) {
	//rootPath, _ := os.Getwd()
	//_ = fmt.Sprintf("%s/%s", rootPath, fileName)
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	for _, env := range config.Envs {
		if err := env.Init(); err != nil {
			return nil, err
		}
	}
	return &config, nil
}

func (c *Config) GetEnvConfig(envType string) (*EnvConfig, error) {
	envConfig, ok := c.Envs[envType]
	if !ok {
		return nil, errors.New("not found env: " + envType)
	}

	if envConfig != nil {
		log.Infoln(fmt.Sprintf("CONFIG envType:%s", envType))
		log.Infoln(fmt.Sprintf("CONFIG APP   AppId:%s,  httpHost:%s,   httpPort:%d,   rootUrl:%s", envConfig.App.AppId, envConfig.App.HttpHost, envConfig.App.HttpPort, envConfig.App.RootUrl))
		log.Infoln(fmt.Sprintf("CONFIG DAPR  host:%s,  httpPort:%d,   grpcPort:%d,   pubsubs:%s",
			envConfig.Dapr.GetHost(), envConfig.Dapr.GetHttpPort(), envConfig.Dapr.GetGrpcPort(), envConfig.Dapr.Pubsubs))
		return envConfig, nil
	}

	return nil, NewEnvTypeError(fmt.Sprintf("error config env-type is \"%s\". choose one of: [dev, test, prod]", envType))
}

type MongoConfig struct {
	Host         string `yaml:"host"`
	Database     string `yaml:"dbname"`
	UserName     string `yaml:"user"`
	Password     string `yaml:"pwd"`
	MaxPoolSize  uint64 `yaml:"maxPoolSize"`
	ReplicaSet   string `yaml:"replicaSet"`
	WriteConcern string `yaml:"writeConcern"`
	ReadConcern  string `yaml:"readConcern"`
}

func (m MongoConfig) IsEmpty() bool {
	if m.Host == "" && m.Database == "" && m.Password == "" && m.UserName == "" {
		return true
	}
	return false
}
