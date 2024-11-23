package rs_server

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/k6/server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/afero"
)

// InitRsServer
//
//	@Description: 添加到restapp的初始化函数 Options.Init
//	@param s
//	@return error
func InitRsServer(httpServer *restapp.HttpServer) error {
	env := httpServer.EnvConfig()
	if env.App.RsServer.Enable {
		fsManger, err := env.GetFsManager()
		if err != nil {
			return err
		}
		fileFs, fileOk := fsManger.Get(env.App.RsServer.FileFsName)
		if !fileOk {
			return errors.New(env.Name + ".app.rsServer fileFsName not found")
		}
		httpFs, httpOk := fsManger.Get(env.App.RsServer.HttpFsName)
		if !httpOk {
			return errors.New(env.Name + ".app.rsServer httpFsName not found")
		}
		if env.App.RsServer.BasePath != "" {
			fileFs = afero.NewBasePathFs(fileFs, env.App.RsServer.BasePath)
		}
		srcFs := NewSrcFsConfig(fileFs, httpFs)
		appEnv := httpServer.EnvConfig()
		runtimeEnv := NewEnv(appEnv)
		data, err := runtimeEnv.ToMap()
		if err != nil {
			return err
		}
		jsServer, err := NewServer(httpServer.App(), data, srcFs, appEnv, "/main.js", appEnv.App.RsServer.Reload)
		if err != nil {
			return err
		}
		if err = jsServer.Run(); err != nil {
			return err
		}
		initHttpServer(httpServer)
	}
	return nil
}

func NewEnv(envCfg common.IEnvConfig) *common.Environment {
	return &common.Environment{
		DaprHost:     envCfg.GetDaprHost(),
		DaprHttpPort: envCfg.GetDaprHttpPort(),
		DaprGrpcPort: envCfg.GetDaprGrpcPort(),
		AppId:        envCfg.GetAppId(),
		AppName:      envCfg.GetAppName(),
		EnvCfg:       envCfg,
	}
}

func initHttpServer(hServer *restapp.HttpServer) {
	if hServer == nil {
		return
	}
	hServer.AddSubscribe(server.GetRegisterSubscribe()...)
	if server.GetServer() != nil && server.GetServer().Swagger() != nil {
		hServer.AddSwagger(server.GetServer().Swagger())
	}
}
