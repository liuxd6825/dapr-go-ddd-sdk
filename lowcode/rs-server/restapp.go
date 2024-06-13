package rs_server

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
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
		fileFs, fileOk := fsManger.Get(env.App.RsServer.FileFsId)
		if !fileOk {
			return errors.New(env.Name + ".app.rsServer fileFsId not found")
		}
		httpFs, httpOk := fsManger.Get(env.App.RsServer.HttpFsId)
		if !httpOk {
			return errors.New(env.Name + ".app.rsServer httpFsId not found")
		}
		if env.App.RsServer.BasePath != "" {
			fileFs = afero.NewBasePathFs(fileFs, env.App.RsServer.BasePath)
		}
		fsCfg := NewFsConfig(fileFs, httpFs)
		envCfg := httpServer.EnvConfig()
		daprClient, err := dapr.GetClient()
		if err != nil {
			return err
		}
		data := map[string]any{
			"DAPR_HOST":      envCfg.Dapr.Host,
			"DAPR_HTTP_PORT": envCfg.Dapr.HttpPort,
			"DAPR_GRPC_PORT": envCfg.Dapr.GrpcPort,
			"APP_ID":         envCfg.App.AppId,
			"APP_NAME":       envCfg.App.AppName,
			"env":            envCfg,
			"daprClient":     daprClient,
		}
		jsServer, err := New(httpServer.App(), data, fsCfg, "/main.js", env.App.RsServer.Reload)
		if err != nil {
			return err
		}
		if err = jsServer.Run(); err != nil {
			return err
		}
		SetSwagger(httpServer)
	}
	return nil
}

func SetSwagger(hServer *restapp.HttpServer) {
	if hServer == nil {
		return
	}
	hServer.AddSubscribe(server.GetRegisterSubscribe()...)
	if server.GetServer() != nil && server.GetServer().Swagger() != nil {
		hServer.AddSwagger(server.GetServer().Swagger())
	}
}
