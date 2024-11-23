package hserver

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	rs_server "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server"
	common2 "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/spf13/afero"
)

// InitServer
//
//	@Description: 添加到restapp的初始化函数 Options.Init
//	@param s
//	@return error
func InitServer(httpServer *restapp.HttpServer) error {
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
		srcFs := NewSrcFs(fileFs, httpFs)
		appEnv := httpServer.EnvConfig()
		runtimeEnv := rs_server.NewEnv(appEnv)
		data, err := runtimeEnv.ToMap()
		if err != nil {
			return err
		}
		_, err = initServer(httpServer.App(), data, srcFs, appEnv, "/src-server/xsrc/server.html", appEnv.App.RsServer.Reload)
		if err != nil {
			return err
		}
	}
	return nil
}

func initServer(app *iris.Application, data map[string]any, srcFs *SrcFs, envCfg common2.IEnvConfig, mainFileName string, reload bool) (*Server, error) {
	server, err := NewServer(app, mainFileName, envCfg, data)
	if err != nil {
		return nil, err
	}
	return server, nil
}
