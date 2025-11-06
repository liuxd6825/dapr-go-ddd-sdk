package restapp

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	dapr2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	logs2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs/userlog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

const SystemTenantId = "system"

func init() {

}
func InitApplication(ctx context.Context, envCfg *EnvConfig, eventTypes []RegisterEventType, isTest bool, fun func(cxt context.Context) error) error {
	if envCfg == nil {
		return errors.New("envConfig is null")
	}

	if len(envCfg.App.HttpHost) == 0 {
		envCfg.App.HttpHost = "0.0.0.0"
	}

	// 解决mongo数据库上日期的时区不一致的问题
	times.SetLocalTimeZone()

	userlog.Init(envCfg.App.AppId, envCfg.App.AppName)

	//设置CPU与内容
	if err := setCpuMemory(envCfg.Name, envCfg.App); err != nil {
		return err
	}

	if envCfg.App.AuthToken != "" {
		DefaultAuthToken = envCfg.App.AuthToken
	}

	e := NewEnv(envCfg)
	e.Init()
	env.SetEnv(e)

	logs2.Infof(ctx, nil, fmt.Sprintf("ctype=app; appId=%s; env=%s;", e.App.AppId, e.Name))
	logs2.Infof(ctx, nil, fmt.Sprintf("ctype=app; httpHost=%s; httpPort=%d; httpRootUrl=%s;", e.App.HttpHost, e.App.HttpPort, e.App.RootUrl))
	logs2.Infof(ctx, nil, fmt.Sprintf("ctype=dapr; daprHost=%s; daprHttpPort=%d; daprGrpcPort=%d;", e.Dapr.GetHost(), e.Dapr.GetHttpPort(), e.Dapr.GetGrpcPort()))
	logs2.Infof(ctx, nil, fmt.Sprintf("ctype=eventStores; length=%v;", len(e.Dapr.EventStores)))

	// 注册领域事件类型
	for _, t := range eventTypes {
		if err := ddd.RegisterEventType(t.EventType, t.Version, t.NewFunc); err != nil {
			return errors.New("RegisterEventType() error:\"%s\" , EventType=\"%s\", Version=\"%s\"", err.Error(), t.EventType, t.Version)
		}
	}

	// 注册事件存储器
	eventStoresMap := newEventStores(e.Dapr, e.Dapr.GetClient())
	for key, es := range eventStoresMap {
		ddd.RegisterEventStore(key, es)
	}
	var err error
	if fun != nil {
		err = fun(ctx)
	}
	return err
}

func initLogs(envConfig *EnvConfig, level logs2.Level, saveDays int, rotationHour int, logFile string, outputType string) error {
	outType, err := logs2.ParseOutputType(outputType)
	if err != nil {
		return err
	}
	logs2.Init(AbsFileName(logFile, envConfig), level, saveDays, rotationHour, outType)
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

	var fields logs2.Fields
	ctx := context.Background()

	if config.CPU != nil {
		cpu, err := setCpu(*config.CPU)
		if err != nil {
			logs2.Errorf(ctx, fields, "ctype=app; envName=%s; cpu=%v; error=%s ", envName, cpu, err.Error())
			return err
		} else {
			logs2.Infof(ctx, fields, "ctype=app; cpu=%v;", cpu)
		}
	}

	if config.Memory != nil {
		memTxt, err := setMem(*config.Memory)
		if err != nil {
			logs2.Errorf(ctx, fields, "ctype=app; memory=%s; error=%s; 值不正确。示例: 10G, 10M, 10K", envName, memTxt, err.Error())
			return err
		} else {
			logs2.Infof(ctx, fields, "ctype=app; memory=%s; ", memTxt)
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

func newEventStores(cfg *env.Dapr, client dapr2.DaprClient) map[string]ddd.EventStore {
	//创建dapr事件存储器
	eventStoresMap := make(map[string]ddd.EventStore)
	if !cfg.IsEnable() {
		return eventStoresMap
	}
	esMap := cfg.EventStores
	if len(esMap) == 0 {
		logs2.Panicf(context.Background(), nil, "config eventStores is empity")
	} else {
		var defEs ddd.EventStore
		for _, item := range esMap {
			eventStorage, err := ddd.NewGrpcEventStore(item.CompName, item.PubSubName, client)
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
