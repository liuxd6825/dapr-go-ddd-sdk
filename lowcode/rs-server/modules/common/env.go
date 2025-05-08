package common

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
)

type Environment struct {
	DaprHost     string
	DaprHttpPort int64
	DaprGrpcPort int64
	AppId        string
	AppName      string
	EnvCfg       IEnvConfig
}

func (e *Environment) ToMap() (map[string]any, error) {
	daprClient, err := dapr.GetClient()
	if err != nil {
		return nil, err
	}
	res := map[string]any{
		"DAPR_HOST":      e.DaprHost,
		"DAPR_HTTP_PORT": e.DaprHttpPort,
		"DAPR_GRPC_PORT": e.DaprGrpcPort,
		"APP_ID":         e.AppId,
		"APP_NAME":       e.AppName,
		"env":            e.EnvCfg,
		"daprClient":     daprClient,
	}
	return res, nil
}
