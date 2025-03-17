package test

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/fsm"
)

type EnvConfig struct {
	AppId           string
	AppName         string
	AppHttpHost     string
	AppHttpPort     int
	DaprHost        string
	DaprHttpPort    int64
	DaprGrpcPort    int64
	RsServerSrcPath string
	RsServerEnable  bool
	Fsm             *fsm.Manager
}

func NewEnvConfig(fsm *fsm.Manager, srcPath string) *EnvConfig {
	return &EnvConfig{
		AppId:           "test",
		AppName:         "test",
		AppHttpHost:     "localhost",
		AppHttpPort:     8080,
		DaprHost:        "localhost",
		DaprHttpPort:    8080,
		DaprGrpcPort:    8080,
		RsServerSrcPath: srcPath,
		RsServerEnable:  true,
		Fsm:             fsm,
	}
}

func (c *EnvConfig) GetAppId() string {
	return c.AppId
}

func (c *EnvConfig) GetAppName() string {
	return c.AppName
}

func (c *EnvConfig) GetAppHttpHost() string {
	return c.AppHttpHost
}

func (c *EnvConfig) GetAppHttpPort() int {
	return c.AppHttpPort
}

func (c *EnvConfig) GetDaprHost() string {
	return c.DaprHost
}

func (c *EnvConfig) GetDaprHttpPort() int64 {
	return c.DaprHttpPort
}

func (c *EnvConfig) GetDaprGrpcPort() int64 {
	return c.DaprGrpcPort
}

func (c *EnvConfig) GetFsManager() (*fsm.Manager, error) {
	return c.Fsm, nil
}

func (c *EnvConfig) GetHServerSrcPath() string {
	return c.RsServerSrcPath
}

func (c *EnvConfig) GetHServerEnable() bool {
	return c.RsServerEnable
}
