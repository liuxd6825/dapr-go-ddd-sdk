package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/processutils"
	"strconv"
)

func start(env *EnvConfig) {
	_ = startDapr(env)
}

func status(env *EnvConfig) {
	fmt.Println("")
	_, _ = statusService(env)

	fmt.Println("")
	_, _ = statusDapr(env)

	fmt.Println("")
}

func stop(env *EnvConfig) error {
	fmt.Println("")
	_ = stopDapr(env)

	fmt.Println("")
	_ = stopService(env)

	fmt.Println("")
	return nil
}

func startDapr(env *EnvConfig) error {
	if env == nil {
		return nil
	}
	if !env.Dapr.Server.Start {
		return nil
	}

	p := newDaprProcess(env)
	err := p.Start()
	if err != nil {
		return err
	}
	return nil
}

func stopDapr(env *EnvConfig) error {
	if env == nil {
		return nil
	}
	p := newDaprProcess(env)
	_, _ = p.Kill()
	return nil
}

func stopService(env *EnvConfig) error {
	if env == nil {
		return nil
	}
	p := newServiceProcess(env)
	_, _ = p.Kill()
	return nil
}

func statusDapr(env *EnvConfig) ([]*processutils.ProcessInfo, error) {
	fmt.Println("Dapr:")
	if env == nil {
		return nil, nil
	}
	p := newDaprProcess(env)
	list, err := p.GetProcessInfo()
	if err != nil {
		fmt.Println("Dapr: 查找Dapr进程时出错, 信息:", err.Error())
	}
	return list, err
}

func statusService(env *EnvConfig) ([]*processutils.ProcessInfo, error) {
	p := newServiceProcess(env)
	fmt.Println("Service:")
	list, err := p.GetProcessInfo()
	if err != nil {
		fmt.Println("Service: 查找服务进程时出错, 信息:", err.Error())
	}
	return list, err
}

func newServiceProcess(env *EnvConfig) processutils.Process {
	exeName := GetExeName()
	p := processutils.NewProcess(exeName, nil, "start")
	return p
}

func newDaprProcess(env *EnvConfig) processutils.Process {
	appId := env.App.AppId
	appHttpPort := strconv.FormatInt(int64(env.App.HttpPort), 10)
	daprHttpPort := strconv.FormatInt(*env.Dapr.HttpPort, 10)
	daprGrpcPort := strconv.FormatInt(*env.Dapr.GrpcPort, 10)
	config := env.Dapr.Server.Config
	componentsPath := env.Dapr.Server.ComponentsPath
	enableMetrics := strconv.FormatBool(env.Dapr.Server.EnableMetrics)
	logLevel := env.Dapr.Server.LogLevel
	placementHostAddress := env.Dapr.Server.PlacementHostAddress
	logFile := env.Dapr.Server.LogFile
	logOutputType := env.Dapr.Server.LogOutputType

	args := []string{
		"-app-id=" + appId,
		"-app-port=" + appHttpPort,
		"-dapr-http-port=" + daprHttpPort,
		"-dapr-grpc-port=" + daprGrpcPort,
		"-log-level=" + logLevel,
		"-log-output-type=" + logOutputType,
		"-log-file=" + AbsFileName(logFile),
		"-enable-metrics=" + enableMetrics,
		"-config=" + AbsFileName(config),
		"-components-path=" + AbsFileName(componentsPath),
		"-placement-host-address=" + placementHostAddress,
	}
	p := processutils.NewProcess("daprd", args, "app-id="+appId, "app-port="+appHttpPort)
	return p
}
