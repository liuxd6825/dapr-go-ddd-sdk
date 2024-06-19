package restapp

type IEnvConfig interface {
	GetAppId() string
	GetAppName() string
	GetAppHttpHost() string
	GetAppHttpPort() int
	GetDaprHost() string
	GetDaprHttpPort() int64
	GetDaprGrpcPort() int64
}

func (e *EnvConfig) RsConfig() IEnvConfig {
	return e
}

func (e *EnvConfig) GetAppId() string {
	return e.App.AppId
}

func (e *EnvConfig) GetAppName() string {
	return e.App.AppName
}

func (e *EnvConfig) GetAppHttpHost() string {
	return e.App.HttpHost
}

func (e *EnvConfig) GetAppHttpPort() int {
	return e.App.HttpPort
}

func (e *EnvConfig) GetDaprHost() string {
	return e.Dapr.GetHost()
}

func (e *EnvConfig) GetDaprHttpPort() int64 {
	return e.Dapr.GetHttpPort()
}

func (e *EnvConfig) GetDaprGrpcPort() int64 {
	return e.Dapr.GetGrpcPort()
}
