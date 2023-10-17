package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
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
	Neo4j map[string]*Neo4jConfig `json:"neo4j"`
	Mysql map[string]*MySqlConfig `json:"mysql"`
	Minio map[string]*MinioConfig `yaml:"minio"`
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

type LogConfig struct {
	Level string `yaml:"level"`
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

func (c *Config) GetEnvConfig(env string) (*EnvConfig, error) {
	envConfig, ok := c.Envs[env]
	if !ok {
		return nil, errors.New("not found env: " + env)
	}

	if envConfig != nil {
		log.Infoln(fmt.Sprintf("config env:%s", env))
		log.Infoln(fmt.Sprintf("config app appId:%s, httpHost:%s, httpPort:%d, rootUrl:%s", envConfig.App.AppId, envConfig.App.HttpHost, envConfig.App.HttpPort, envConfig.App.RootUrl))
		log.Infoln(fmt.Sprintf("config dapr host:%s,  httpPort:%d, grpcPort:%d, pubsubs:%s",
			envConfig.Dapr.GetHost(), envConfig.Dapr.GetHttpPort(), envConfig.Dapr.GetGrpcPort(), envConfig.Dapr.Pubsubs))
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
