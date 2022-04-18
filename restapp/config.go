package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type Config struct {
	EnvName string    `yaml:"env"`
	Test    EnvConfig `yaml:"test"`
	Dev     EnvConfig `yaml:"dev"`
	Prod    EnvConfig `yaml:"prod"`
}

type EnvConfig struct {
	App   AppConfig   `yaml:"app"`
	Log   LogConfig   `yaml:"log"`
	Dapr  DaprConfig  `yaml:"dapr"`
	Mongo MongoConfig `yaml:"mongo"`
}

func (e EnvConfig) CheckError() error {
	if e.Log.Level != "" {
		l, err := applog.NewLevel(e.Log.Level)
		if err != nil {
			return err
		}
		e.Log.level = l
	}
	return nil
}

type AppConfig struct {
	AppId   string `yaml:"id"`
	AppPort int    `yaml:"http-port"`
	RootUrl string `yaml:"root-url"`
}

type DaprConfig struct {
	Host     string   `yaml:"host"`
	HttpPort int      `yaml:"http-port"`
	GrpcPort int      `yaml:"grpc-port"`
	Pubsubs  []string `yaml:"pubsubs,flow"`
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

	if err := config.Dev.CheckError(); err != nil {
		return nil, err
	}
	if err := config.Test.CheckError(); err != nil {
		return nil, err
	}
	if err := config.Prod.CheckError(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) GetEnvConfig() (*EnvConfig, error) {
	envName := strings.ToLower(c.EnvName)
	switch envName {
	case "test":
		return &c.Test, nil
	case "dev":
		return &c.Dev, nil
	case "prod":
		return &c.Prod, nil
	}
	return nil, NewEnvNameError(fmt.Sprintf("config.envName is \"%s\" error. range is [dev, test, prod]", c.EnvName))
}

type MongoConfig struct {
	Host        string `yaml:"host"`
	Database    string `yaml:"dbname"`
	UserName    string `yaml:"user"`
	Password    string `yaml:"pwd"`
	MaxPoolSize uint64 `yaml:"max-pool-size"`
}

func (m MongoConfig) IsEmpty() bool {
	if m.Host == "" && m.Database == "" && m.Password == "" && m.UserName == "" {
		return true
	}
	return false
}
