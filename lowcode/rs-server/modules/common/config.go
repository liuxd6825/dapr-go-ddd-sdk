package common

type IEnvConfig interface {
	GetAppId() string
	GetAppName() string
	GetAppHttpHost() string
	GetAppHttpPort() int
	GetDaprHost() string
	GetDaprHttpPort() int64
	GetDaprGrpcPort() int64
}
