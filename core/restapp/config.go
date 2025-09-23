package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Env  string                `yaml:"env"`
	Envs map[string]*EnvConfig `yaml:"envs"`
}

type EnvConfig struct {
	Name      string                     `yaml:"-" json:"name"`
	App       *AppConfig                 `yaml:"app" json:"app"`
	Log       *LogConfig                 `yaml:"log" json:"log"`
	Dapr      *DaprConfig                `yaml:"dapr" json:"dapr"`
	Resources map[string]*ResourceConfig `yaml:"resources" json:"resources"`
	Mongo     map[string]*MongoConfig    `yaml:"mongo" json:"mongo"`
	Neo4j     map[string]*Neo4jConfig    `yaml:"neo4j" json:"neo4J"`
	Mysql     map[string]*MySqlConfig    `yaml:"mysql" json:"mysql"`
	Minio     map[string]*MinioConfig    `yaml:"minio" json:"minio"`
	Redis     map[string]*RedisConfig    `yaml:"redis" json:"redis"`
	Fs        []map[string]any           `yaml:"fs" json:"fs"`
	Auth      *AuthConfig                `yaml:"auth" json:"auth"`
	Temporal  *TemporalConfig            `yaml:"temporal" json:"temporal"`
}

func NewEnvConfig(name string) *EnvConfig {
	return &EnvConfig{
		Name:      name,
		App:       NewAppConfig(),
		Log:       newLogConfig(),
		Dapr:      newDaprConfig(),
		Resources: map[string]*ResourceConfig{},
		Mongo:     map[string]*MongoConfig{},
		Neo4j:     map[string]*Neo4jConfig{},
		Mysql:     map[string]*MySqlConfig{},
		Minio:     map[string]*MinioConfig{},
		Redis:     map[string]*RedisConfig{},
		Fs:        []map[string]any{},
		Auth:      &AuthConfig{},
		Temporal:  &TemporalConfig{},
	}
}

func newLogConfig() *LogConfig {
	return &LogConfig{
		Level:      "info",
		SaveDays:   10,
		SplitHour:  12,
		OutputType: "console",
		level:      logs.InfoLevel,
	}
}

func newDaprConfig() *DaprConfig {
	enable := false
	return &DaprConfig{
		Enable: &enable,
	}
}

// FsRootPath fs文件系统取根路径
type FsRootPath interface {
	GetRootPath() string
}

// AppConfig
// @Description:  应用配置
// @Author:       liuxdl
// @Date:         2021/10/18 10:57
type AppConfig struct {
	AppId      string         `yaml:"id" json:"id"`                 // 应用ID
	AppName    string         `yaml:"name" json:"name"`             // 应用名称
	ProdMode   bool           `yaml:"prodMode" json:"prodMode"`     // 是生产模式
	HttpHost   string         `yaml:"httpHost" json:"httpHost"`     // 绑定HTTP IP
	HttpPort   int            `yaml:"httpPort" json:"httpPort"`     // 绑定HTTP 端口
	RootUrl    string         `yaml:"rootUrl" json:"rootUrl"`       // URL根
	CPU        *int           `yaml:"cpu" json:"cpu"`               // CPU数量
	Memory     *string        `yaml:"memory" json:"memory"`         // 内存大小
	Meta       map[string]any `yaml:"meta" json:"meta"`             // 自定义配置
	AuthToken  string         `yaml:"authToken" json:"authToken"`   // 开发时Token
	HServer    HServer        `yaml:"hServer" json:"hServer"`       // 脚本服务配置
	Template   HtmlTemplate   `yaml:"template" json:"template"`     // html模板配置
	IsPubEvent bool           `yaml:"isPubEvent" json:"isPubEvent"` // 是否发送领域事件
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		AppId:    "",
		AppName:  "",
		HttpHost: "0.0.0.0",
		HttpPort: 1984,
	}
}

// HServer
// @Description: 脚本服务配置
// @Author:       liuxd
// @Date:         2021/10/18 10:57
type HServer struct {
	Enable       bool           `yaml:"enable" json:"enable"`             // 是否启用脚本服务
	SrcName      string         `yaml:"srcName" json:"srcName"`           // API源码文件系统名称
	WebName      string         `yaml:"webName" json:"webName"`           // Web源码文件系统名称
	BasePath     string         `yaml:"basePath" json:"basePath"`         // 脚本文件路径
	Reload       bool           `yaml:"reload" json:"reload"`             // 是否自动加载脚本
	WatchRestart bool           `yaml:"watchRestart" json:"watchRestart"` // 检查文件变化，重新启动
	Meta         map[string]any `yaml:"meta" json:"meta"`                 // 扩展信息
	Npm          Npm            `yaml:"npm" json:"npm"`                   // 前端npm环境配置
}
type Npm struct {
	Links []*NpmLink `yaml:"links" json:"links"`
}

type NpmLink struct {
	Name string `yaml:"name" json:"name"`
	Path string `yaml:"path" json:"path"`
}

type IHServer interface {
	GetEnable() bool
	GetBasePath() string
}

// HtmlTemplate
// @Description: html 模板配置
// @Author:       liuxd
// @Date:         2021/10/18 10:57
type HtmlTemplate struct {
	Enable   bool   `yaml:"enable" json:"enable"`     // 是否启用html模板
	ApiUrl   string `yaml:"apiUrl" json:"apiUrl"`     // api html模板文件路径
	FsId     string `yaml:"fsId" json:"fsId"`         // 在fs节中配置key
	BasePath string `yaml:"basePath" json:"basePath"` // 脚本文件路径
}

type ResourceConfig struct {
	Namespace string            `yaml:"namespace" json:"namespace"`
	Name      string            `yaml:"name" json:"name"`
	Type      string            `yaml:"type" json:"type"`
	URI       string            `yaml:"uri" json:"uri"`
	Metadata  map[string]string `yaml:"metadata" json:"metadata"`
}

type Metadata map[string]string
type DaprConfig struct {
	Host                *string                `yaml:"host" json:"host"`
	HttpPort            *int64                 `yaml:"httpPort" json:"httpPort"`
	GrpcPort            *int64                 `yaml:"grpcPort" json:"grpcPort"`
	MaxCallRecvMsgSize  *int                   `yaml:"maxCallRecvMsgSize" json:"maxCallRecvMsgSize"` //dapr数据包大小，单位M
	MaxIdleConns        *int                   `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxIdleConnsPerHost *int                   `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost"`
	IdleConnTimeout     *int                   `yaml:"idleConnTimeout" json:"idleConnTimeout"`
	EventStores         map[string]*EventStore `yaml:"eventStores" json:"eventStores"`
	Actor               ActorConfig            `yaml:"actor" json:"actor"`
	Start               bool                   `yaml:"start" json:"start"`
	StartArgs           Metadata               `yaml:"startArgs" json:"startArgs"`
	Enable              *bool                  `yaml:"enable" json:"enable"`
}

// DaprServerConfig dapr服务端参数
type DaprServerConfig struct {
	Start                bool   `yaml:"start" json:"start"`                               //是否启动Dapr
	EnableMetrics        bool   `yaml:"enableMetrics" json:"enableMetrics"`               //是否健康检测
	Config               string `yaml:"config" json:"config"`                             // 配置文件参数
	ComponentsPath       string `yaml:"componentsPath" json:"componentsPath"`             // 组件参数
	PlacementHostAddress string `yaml:"placementHostAddress" json:"placementHostAddress"` // placement地址
	LogLevel             string `yaml:"logLevel" json:"logLevel"`                         // 日志级别
	LogFile              string `yaml:"logFile" json:"logFile"`                           // 日志文件
	LogOutputType        string `yaml:"logOutputType" json:"logOutputType"`               // 输出日期类型
}

// ActorConfig dapr actor配置
type ActorConfig struct {
	ActorIdleTimeout       string `yaml:"actorIdleTimeout" json:"actorIdleTimeout"`
	ActorScanInterval      string `yaml:"actorScanInterval" json:"actorScanInterval"`
	DrainOngingCallTimeout string `yaml:"drainOngoingCallTimeout" json:"drainOngingCallTimeout"`
	DrainBalancedActors    bool   `yaml:"drainRebalancedActors" json:"drainBalancedActors"`
}

// EventStore 事件存储
type EventStore struct {
	CompName   string `yaml:"name" json:"name"`     // Dapr EventStarge 组件名称
	PubSubName string `yaml:"pubsub" json:"pubsub"` // Dapr Pubsub 组件名称
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" json:"level"`
	SaveDays   int    `yaml:"saveDays" json:"saveDays"`   //日志保存的天数
	SplitHour  int    `yaml:"splitHour" json:"splitHour"` //文件分隔时间，单位小时
	LogFile    string `yaml:"logFile" json:"logFile"`
	OutputType string `yaml:"outputType" json:"outputType"` // 日志输出类型 console、 file、 all
	level      logs.Level
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
	// 转换为绝对路径
	filename, err := filepath.Abs(filename)
	yamlFile, err := os.ReadFile(filename)
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
	for _, m := range e.Fs {
		for k, v := range m {
			if s, ok := v.(string); ok {
				m[k] = ReplaceSysValues(s)
			}
		}
	}
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
	if e.Log.LogFile == "" {
		e.Log.LogFile = fmt.Sprintf("./logs/%s.log", GetExeName())
	}
	e.Log.level = level

	//初始化Dapr
	/*if err := e.Dapr.init(e); err != nil {
		return err
	}*/

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

/*
	func (e *EnvConfig) GetEnvString(envName string, defValue *string) *string {
		value, ok := os.LookupEnv(envName)
		if !ok {
			return defValue
		}
		return &value
	}

	func (e *EnvConfig) GetProdMode() bool {
		return e.App.ProdMode
	}

	func (e *EnvConfig) GetMeta() map[string]any {
		return e.App.Meta
	}

	func (s *HServer) GetEnable() bool {
		return s.Enable
	}

	func (s *HServer) GetBasePath() string {
		return s.BasePath
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
			val := dapr2.GetMaxCallRecvMsgSize()
			c.MaxCallRecvMsgSize = &val
		}

		if c.MaxIdleConnsPerHost == nil {
			val := dapr2.DefaultMaxIdleConnsPerHost
			c.MaxIdleConns = &val
		}

		if c.IdleConnTimeout == nil {
			val := dapr2.DefaultIdleConnTimeout
			c.IdleConnTimeout = &val
		}

		if c.MaxIdleConns == nil {
			val := dapr2.DefaultMaxIdleConns
			c.MaxIdleConnsPerHost = &val
		}

		if len(c.EventStores) > 0 {
			for compName, es := range e.Dapr.EventStores {
				if es.CompName == "" {
					es.CompName = compName
				}
				if len(es.PubSubName) == 0 {
					return errors.ErrorOf("config env:%s  Dapr.EventStores.%s pubsub is null", e.Name, compName)
				}
			}
		}

		e.Dapr.Actor.init()

		return nil

}

	func (c *DaprConfig) IsEnable() bool {
		return c.Enable != nil && *c.Enable
	}

	func (c *DaprConfig) GetHost() string {
		if c.Host == nil {
			return ""
		}
		return *c.Host
	}

	func (c *DaprConfig) GetHttpPort() int64 {
		if c.HttpPort == nil {
			return 0
		}
		return *c.HttpPort
	}

	func (c *DaprConfig) GetGrpcPort() int64 {
		if c.GrpcPort == nil {
			return 0
		}
		return *c.GrpcPort
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
*/
func searchConfigFile(path, configName string, fileName string) (string, bool, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", false, err
	}
	if len(files) <= 0 {
		return "", false, nil
	}
	for _, file := range files {
		name := file.Name()
		if file.IsDir() && name == configName {
			list, err := os.ReadDir(path + "/" + file.Name())
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
