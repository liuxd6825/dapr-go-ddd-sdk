package common

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/fsm"
)

type IEnvConfig interface {
	GetAppId() string
	GetAppName() string

	GetAppHttpHost() string
	GetAppHttpPort() int
	GetDaprHost() string
	GetDaprHttpPort() int64
	GetDaprGrpcPort() int64
	GetFsManager() *fsm.Manager
	GetHServerSrcPath() string
	GetHServerEnable() bool
}
