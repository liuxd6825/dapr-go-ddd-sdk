package restapp

type IEnvConfig interface {
	GetAppId() string
	GetAppName() string
	GetAppHttpHost() string
	GetAppHttpPort() int
	GetDaprHost() string
	GetDaprHttpPort() int64
	GetDaprGrpcPort() int64
	GetHServerSrcPath() string
	GetHServerEnable() bool
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

/*
func (e *EnvConfig) GetHServerSrcPath() string {
	return e.App.HServer.GetBasePath()
}

func (e *EnvConfig) GetHServerEnable() bool {
	return e.App.HServer.GetEnable()
}
*/

func (e *EnvConfig) GetHServerSrcPath() string {
	return e.App.HServer.GetBasePath()
	/*
		fsName := e.App.HServer.SrcName
		fsVal, ok := e.fsManager.GetFs(fsName)
		if !ok {
			return ""
		}
		if path, ok := fsVal.(FsRootPath); ok {
			return path.GetRootPath()
		} else {
			panic(fmt.Sprintf("fs \"%s\" does not have a GetRootPath() method", fsName))
		}
		return "" */
}

func (e *EnvConfig) GetHServerEnable() bool {
	return e.App.HServer.Enable
}
