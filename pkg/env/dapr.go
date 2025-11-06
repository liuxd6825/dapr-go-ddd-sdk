package env

import (
	"context"
	"fmt"
	"os"
	"strconv"

	dapr2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/intutils"
)

type Dapr struct {
	Host string `yaml:"host" json:"host"`

	HttpPort            int64                      `yaml:"httpPort" json:"httpPort"`
	GrpcPort            int64                      `yaml:"grpcPort" json:"grpcPort"`
	MaxCallRecvMsgSize  int                        `yaml:"maxCallRecvMsgSize" json:"maxCallRecvMsgSize"` //dapr数据包大小，单位M
	MaxIdleConns        int                        `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxIdleConnsPerHost int                        `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost"`
	IdleConnTimeout     int                        `yaml:"idleConnTimeout" json:"idleConnTimeout"`
	EventStores         map[string]*DaprEventStore `yaml:"eventStores" json:"eventStores"`
	Actor               DaprActor                  `yaml:"actor" json:"actor"`
	Start               bool                       `yaml:"start" json:"start"`
	StartArgs           map[string]any             `yaml:"startArgs" json:"startArgs"`
	//Enable              bool                       `yaml:"enable" json:"enable"`
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	client     dapr2.DaprClient
}

// DaprServer dapr服务端参数
type DaprServer struct {
	Start                bool   `yaml:"start" json:"start"`                               //是否启动Dapr
	EnableMetrics        bool   `yaml:"enableMetrics" json:"enableMetrics"`               //是否健康检测
	Config               string `yaml:"config" json:"config"`                             // 配置文件参数
	ComponentsPath       string `yaml:"componentsPath" json:"componentsPath"`             // 组件参数
	PlacementHostAddress string `yaml:"placementHostAddress" json:"placementHostAddress"` // placement地址
	LogLevel             string `yaml:"logLevel" json:"logLevel"`                         // 日志级别
	LogFile              string `yaml:"logFile" json:"logFile"`                           // 日志文件
	LogOutputType        string `yaml:"logOutputType" json:"logOutputType"`               // 输出日期类型
}

// DaprActor dapr actor配置
type DaprActor struct {
	ActorIdleTimeout       string `yaml:"actorIdleTimeout" json:"actorIdleTimeout"`
	ActorScanInterval      string `yaml:"actorScanInterval" json:"actorScanInterval"`
	DrainOngingCallTimeout string `yaml:"drainOngoingCallTimeout" json:"drainOngingCallTimeout"`
	DrainBalancedActors    bool   `yaml:"drainRebalancedActors" json:"drainBalancedActors"`
}

// DaprEventStore 事件存储
type DaprEventStore struct {
	CompName   string `yaml:"name" json:"name"`     // Dapr EventStarge 组件名称
	PubSubName string `yaml:"pubsub" json:"pubsub"` // Dapr Pubsub 组件名称
}

func initDapr(e *Env) {
	if e.Dapr == nil {
		e.Dapr = NewDapr()
		return
	}
	c := e.Dapr
	if c.Host == "" {
		var value = "localhost"
		c.Host = getEnvString("DAPR_HOST", value)
	}
	if c.ApiVersion == "" {
		c.ApiVersion = "v1.0"
	}

	if e.Dapr.HttpPort < 0 {
		var value int64 = 3500
		c.HttpPort = getEnvInt("DAPR_HTTP_PORT", value)
	}

	if e.Dapr.GrpcPort < 0 {
		var value int64 = 50001
		c.GrpcPort = getEnvInt("DAPR_GRPC_PORT", value)
	}

	if c.MaxCallRecvMsgSize <= 0 {
		val := dapr2.GetMaxCallRecvMsgSize()
		c.MaxCallRecvMsgSize = val
	}

	if c.MaxIdleConnsPerHost <= 0 {
		val := dapr2.DefaultMaxIdleConnsPerHost
		c.MaxIdleConns = val
	}

	if c.IdleConnTimeout <= 0 {
		val := dapr2.DefaultIdleConnTimeout
		c.IdleConnTimeout = val
	}

	if c.MaxIdleConns <= 0 {
		val := dapr2.DefaultMaxIdleConns
		c.MaxIdleConnsPerHost = val
	}

	if len(c.EventStores) > 0 {
		for compName, es := range e.Dapr.EventStores {
			if es.CompName == "" {
				es.CompName = compName
			}
			if len(es.PubSubName) == 0 {
				err := errors.ErrorOf("config env:%s  Dapr.EventStores.%s pubsub is null", e.Name, compName)
				panic(err)
			}
		}
	}

	e.Dapr.Actor.init()

	var client dapr2.DaprClient
	var err error
	var ctx = context.Background()
	if e.Dapr.IsEnable() {
		// 启动服务，创建dapr客户端
		client, err = dapr2.NewDaprClient(ctx, e.Dapr.GetHost(), e.Dapr.GetHttpPort(), e.Dapr.GetGrpcPort(), func(ops *dapr2.DaprHttpOptions) {
			ops.MaxCallRecvMsgSize = intutils.P2IntDefault(&e.Dapr.MaxCallRecvMsgSize, dapr2.GetMaxCallRecvMsgSize())
			ops.MaxIdleConns = intutils.P2IntDefault(&e.Dapr.MaxIdleConns, dapr2.DefaultMaxIdleConns)
			ops.MaxIdleConnsPerHost = intutils.P2IntDefault(&e.Dapr.MaxIdleConnsPerHost, dapr2.DefaultMaxIdleConnsPerHost)
			ops.IdleConnTimeout = intutils.P2IntDefault(&e.Dapr.IdleConnTimeout, dapr2.DefaultIdleConnTimeout)
		})
		if err != nil {
			panic(err)
		}
		dapr2.SetDaprClient(client)
		e.Dapr.client = client
	}
}

func (c *Dapr) IsEnable() bool {
	return c.Start
}

func (c *Dapr) GetHost() string {
	return c.Host
}

func (c *Dapr) GetHttpPort() int64 {
	return c.HttpPort
}

func (c *Dapr) GetGrpcPort() int64 {
	return c.GrpcPort
}

func (c *Dapr) GetClient() dapr2.DaprClient {
	return c.client
}

func (c *DaprActor) init() {
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

func (c *Dapr) GetInvokeService(serviceName string, method string) string {
	return fmt.Sprintf("http://%s:%d/v1.0/invoke/%s/method/%s", c.Host, c.HttpPort, serviceName, method)
}

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

func getEnvInt(envName string, defValue int64) int64 {
	value, ok := os.LookupEnv(envName)
	if !ok {
		return defValue
	}
	parseInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}
	return parseInt
}

func getEnvString(envName string, defValue string) string {
	value, ok := os.LookupEnv(envName)
	if !ok {
		return defValue
	}
	return value
}
