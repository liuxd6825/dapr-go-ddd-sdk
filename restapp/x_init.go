package restapp

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs/userlog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/setting"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	_envConfig *EnvConfig
)

const SystemTenantId = "system"

func InitApplication(ctx context.Context, envConfig *EnvConfig, eventTypes []RegisterEventType, isTest bool, fun func(cxt context.Context) error) error {
	if envConfig == nil {
		return errors.New("envConfig is null")
	}

	if len(envConfig.App.HttpHost) == 0 {
		envConfig.App.HttpHost = "0.0.0.0"
	}

	// 设置全局时区为本地时区
	setting.SetLocalTimeZone()

	userlog.Init(envConfig.App.AppId, envConfig.App.AppName)

	//设置CPU与内容
	if err := setCpuMemory(envConfig.Name, &envConfig.App); err != nil {
		return err
	}

	if err := initMongo(envConfig.App.AppId, envConfig.Mongo); err != nil {
		return err
	}

	if err := initNeo4j(envConfig.Neo4j); err != nil {
		return err
	}

	if err := initResources(envConfig.Resources); err != nil {
		return err
	}

	if err := initMinio(envConfig.Minio); err != nil {
		return err
	}

	if err := initRedis(envConfig.Redis); err != nil {
		return err
	}

	if envConfig.App.AuthToken != "" {
		DefaultAuthToken = envConfig.App.AuthToken
	}

	SetEnvConfig(envConfig)

	// 启动服务，创建dapr客户端
	daprClient, err := dapr.NewDaprClient(ctx, envConfig.Dapr.GetHost(), envConfig.Dapr.GetHttpPort(), envConfig.Dapr.GetGrpcPort(), func(ops *dapr.DaprHttpOptions) {
		ops.MaxCallRecvMsgSize = intutils.P2IntDefault(envConfig.Dapr.MaxCallRecvMsgSize, dapr.GetMaxCallRecvMsgSize())
		ops.MaxIdleConns = intutils.P2IntDefault(envConfig.Dapr.MaxIdleConns, dapr.DefaultMaxIdleConns)
		ops.MaxIdleConnsPerHost = intutils.P2IntDefault(envConfig.Dapr.MaxIdleConnsPerHost, dapr.DefaultMaxIdleConnsPerHost)
		ops.IdleConnTimeout = intutils.P2IntDefault(envConfig.Dapr.IdleConnTimeout, dapr.DefaultIdleConnTimeout)
	})

	if err != nil {
		return err
	}

	dapr.SetDaprClient(daprClient)
	ddd.Init(envConfig.App.AppId)
	applog.Init(daprClient, envConfig.App.AppId, envConfig.Log.level)

	// 注册领域事件类型
	for _, t := range eventTypes {
		if err := ddd.RegisterEventType(t.EventType, t.Version, t.NewFunc); err != nil {
			return errors.New(fmt.Sprintf("RegisterEventType() error:\"%s\" , EventType=\"%s\", Version=\"%s\"", err.Error(), t.EventType, t.Version))
		}
	}

	// 注册事件存储器
	eventStoresMap := newEventStores(&envConfig.Dapr, daprClient)
	for key, es := range eventStoresMap {
		ddd.RegisterEventStore(key, es)
	}

	if fun != nil {
		err = fun(ctx)
	}
	return err
}

func initLogs(level logs.Level, saveDays int, rotationHour int) error {
	appPath, err := os.Executable()
	if err != nil {
		return err
	}
	//appPath = filepath.Clean(appPath)
	appPath, err = filepath.Abs(appPath)
	if err != nil {
		return err
	}

	dir, appName := filepath.Split(appPath)
	saveFile := fmt.Sprintf("%s/logs/%s", dir, appName)
	logs.Init(saveFile, level, saveDays, rotationHour)
	return nil
}

// setCpuMemory
//
//	@Description: 设置Cpu和内存大小
//	@param config
func setCpuMemory(envName string, config *AppConfig) error {
	if config == nil || (config.CPU == nil && config.Memory == nil) {
		// set GOMAXPROCS
		// 适用docker环境
		//_, _ = maxprocs.Set()
		return nil
	}

	var fields logs.Fields
	ctx := context.Background()

	if config.CPU != nil {
		cpu, err := setCpu(*config.CPU)
		if err != nil {
			logs.Panic(ctx, "", fields, "ctype=app; cpu=%v; error=%s ", envName, cpu, err.Error())
			return err
		} else {
			logs.Infof(ctx, "", fields, "ctype=app; cpu=%v;", cpu)
		}
	}

	if config.Memory != nil {
		memTxt, err := setMem(*config.Memory)
		if err != nil {
			logs.Panic(ctx, "", fields, "ctype=app; memory=%s; error=%s; 值不正确。示例: 10G, 10M, 10K", envName, memTxt, err.Error())
			return err
		} else {
			logs.Info(ctx, "", fields, "ctype=app; memory=%s; ", envName, memTxt)
		}
	}

	return nil

}

func setCpu(cpu int) (int, error) {
	maxCpu := runtime.NumCPU()
	if cpu < 0 {
		cpu = maxCpu - cpu
	}
	if cpu > maxCpu {
		cpu = maxCpu
	}
	if cpu <= 0 {
		cpu = 1
	}
	runtime.GOMAXPROCS(cpu)
	return cpu, nil
}

func setMem(val string) (string, error) {
	memTxt := strings.ToLower(strings.Trim(val, " "))
	if memTxt == "" {
		return "", nil
	}
	var memSize int64 = 0
	size := len(memTxt)
	unit := memTxt[size-1 : size]
	memVal := memTxt[0 : size-1]
	memSize, err := stringutils.ToInt64(memVal)
	if err != nil {
		return "", err
	}

	switch unit {
	case "g":
		memSize = memSize * 1024 * 1024 * 1024
	case "m":
		memSize = memSize * 1024 * 1024
	case "k":
		memSize = memSize * 1024
	default:
		return "", errors.New("格式不正确。示例: 10G, 10M, 10K")
	}
	debug.SetMemoryLimit(memSize)
	return memTxt, nil
}

func newEventStores(cfg *DaprConfig, client dapr.DaprClient) map[string]ddd.EventStore {
	//创建dapr事件存储器
	eventStoresMap := make(map[string]ddd.EventStore)
	esMap := cfg.EventStores
	if len(esMap) == 0 {
		logs.Panicf(context.Background(), "", nil, "config eventStores is empity")
	} else {
		var defEs ddd.EventStore
		for _, item := range esMap {
			eventStorage, err := ddd.NewGrpcEventStore(item.CompName, item.PubsubName, client)
			if err != nil {
				panic(err)
			}
			eventStoresMap[item.CompName] = eventStorage
			if defEs == nil {
				defEs = eventStorage
			}
		}
		eventStoresMap[""] = defEs
	}
	return eventStoresMap
}

func GetAppValue(name string) (string, error) {
	var err error
	v, ok := _envConfig.App.Values[name]
	if !ok {
		err = errors.New(fmt.Sprintf("配置变量%s不存在", name))
	}
	return v, err
}

func GetAppValues() map[string]string {
	return _envConfig.App.Values
}

func GetEnvConfig() *EnvConfig {
	return _envConfig
}

func SetEnvConfig(envConfig *EnvConfig) *EnvConfig {
	return _envConfig
}

func GetDaprHost() string {
	return _envConfig.Dapr.GetHost()
}

func GetDaprHttpPort() int64 {
	return _envConfig.Dapr.GetHttpPort()
}

func GetDaprGrpcPort() int64 {
	return _envConfig.Dapr.GetGrpcPort()
}

func GetAppId() string {
	return _envConfig.App.AppId
}

func GetAppHttpHost() string {
	return _envConfig.App.HttpHost
}

func GetHttpInvoke(appId string) string {
	return fmt.Sprintf("http://%s:%v/v1.0/invoke/%v/method/", GetDaprHost(), GetDaprHttpPort(), appId)
}

func GetHttpsInvoke(appId string) string {
	return fmt.Sprintf("https://%s:%v/v1.0/invoke/%v/method/", GetDaprHost(), GetDaprHttpPort(), appId)
}

func GetLogger() logs.Logger {
	return logs.GetLogger()
}
