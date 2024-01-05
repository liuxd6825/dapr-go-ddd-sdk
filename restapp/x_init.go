package restapp

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/daprclient"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/setting"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
)

func init() {

}

func initLogs(level logs.Level, saveDays int, rotationHour int) {
	appPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir, appName := filepath.Split(appPath)
	saveFile := fmt.Sprintf("%s/logs/%s", dir, appName)
	logs.Init(saveFile, level, saveDays, rotationHour)
}

func GetLogger() logs.Logger {
	return logs.GetLogger()
}

func InitApplication(ctx context.Context, envConfig *EnvConfig, eventTypes []RegisterEventType, fun func(cxt context.Context) error) error {
	if envConfig == nil {
		return errors.New("envConfig is null")
	}

	if len(envConfig.App.HttpHost) == 0 {
		envConfig.App.HttpHost = "0.0.0.0"
	}
	// 设置全局时区为本地时区
	setting.SetLocalTimeZone()

	setCpuMemory(envConfig.Name, &envConfig.App)

	if len(envConfig.Mongo) > 0 {
		initMongo(envConfig.App.AppId, envConfig.Mongo)
	}

	if len(envConfig.Neo4j) > 0 {
		initNeo4j(envConfig.Neo4j)
	}

	if len(envConfig.Minio) > 0 {
		if err := initMinio(envConfig.Minio); err != nil {
			return err
		}
	}
	if len(envConfig.Redis) > 0 {
		if err := initRedis(envConfig.Redis); err != nil {
			return err
		}
	}

	if envConfig.App.AuthToken != "" {
		DefaultAuthToken = envConfig.App.AuthToken
	}

	SetCurrentEnvConfig(envConfig)
	
	// 启动服务，创建dapr客户端
	daprClient, err := daprclient.NewDaprDddClient(ctx, envConfig.Dapr.GetHost(), envConfig.Dapr.GetHttpPort(), envConfig.Dapr.GetGrpcPort())
	if err != nil {
		return err
	}

	daprclient.SetDaprDddClient(daprClient)
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

// setCpuMemory
//
//	@Description: 设置Cpu和内存大小
//	@param config
func setCpuMemory(envName string, config *AppConfig) {
	if config == nil {
		return
	}
	var fields logs.Fields
	cpu := config.CPU
	maxCpu := runtime.NumCPU()
	if cpu < 0 {
		cpu = maxCpu - cpu
	}
	if maxCpu < cpu {
		cpu = maxCpu
	}
	if cpu <= 0 {
		cpu = 1
	}
	runtime.GOMAXPROCS(cpu)

	memTxt := strings.ToLower(strings.Trim(config.Memory, " "))
	if memTxt == "" {
		logs.Infof(context.Background(), "", fields, "ctype=app; cpu=%v;", cpu)
		return
	}
	var memSize int64 = 0
	size := len(memTxt)
	unit := memTxt[size-1 : size]
	memVal := memTxt[0 : size-1]
	memSize, err := stringutils.ToInt64(memVal)
	if err != nil {
		logs.Panic(context.Background(), "", fields, "ctype=app; memory=%s; 值不正确。示例: 10G, 10M, 10K", envName, memTxt)
	}

	switch unit {
	case "g":
		memSize = memSize * 1024 * 1024 * 1024
	case "m":
		memSize = memSize * 1024 * 1024
	case "k":
		memSize = memSize * 1024
	default:
		logs.Panic(context.Background(), "", fields, "ctype=app; %s.app.memory=%s 不正确。示例: 10G, 10M, 10K", envName, memTxt)
	}
	debug.SetMemoryLimit(memSize)
	logs.Infof(context.Background(), "", fields, "ctype=app; cpu=%v; memory=%s;", cpu, memTxt)
}
