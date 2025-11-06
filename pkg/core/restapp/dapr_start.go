package restapp

import (
	"fmt"

	"strconv"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/processutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
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
	if !env.Dapr.Start {
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

	/*
		config := AbsFileName(env.Dapr.Server.Config)
		componentsPath := AbsFileName(env.Dapr.Server.ComponentsPath)
		logFile := AbsFileName(env.Dapr.Server.LogFile)
		enableMetrics := strconv.FormatBool(env.Dapr.Server.EnableMetrics)
		logLevel := env.Dapr.Server.LogLevel
		placementHostAddress := env.Dapr.Server.PlacementHostAddress
		logOutputType := env.Dapr.Server.LogOutputType
	*/

	args := []string{
		"-app-id=" + appId,
		"-app-port=" + appHttpPort,
		"-dapr-http-port=" + daprHttpPort,
		"-dapr-grpc-port=" + daprGrpcPort,
		/*
			"-log-level=" + logLevel,
			"-log-output-type=" + logOutputType,
			"-log-file=" + logFile,
			"-enable-metrics=" + enableMetrics,
			"-config=" + config,
			"-components-path=" + componentsPath,
			"-placement-host-address=" + placementHostAddress,

		*/
	}
	if env.Dapr.StartArgs != nil {
		for k, v := range env.Dapr.StartArgs {
			arg := getDaprArg(k, v, env)
			args = append(args, arg)
		}
	}

	//ctx := context.Background()

	//line := strings.Join(args, " ")
	//logs.InfoMsg(ctx, "", "daprd "+line)
	p := processutils.NewProcess("daprd", args, "app-id="+appId, "app-port="+appHttpPort)
	return p
}

func getDaprArg(argName string, argValue string, env *EnvConfig) string {
	argName = stringutils.ToKebabCase(argName)
	if strings.HasPrefix(argValue, "./") || strings.HasPrefix(argValue, "../") {
		argValue = AbsFileName(argValue, env)
	}
	argValue = ReplaceSysValues(argValue, env)
	return fmt.Sprintf("-%s=%s", argName, argValue)
}
