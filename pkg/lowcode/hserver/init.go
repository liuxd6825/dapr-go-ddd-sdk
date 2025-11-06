package hserver

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/handler/file_handler"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type InitOptions func(server element.Server) error

// InitHServer
//
//	@Description:
//	@param fileName
//	@param srcFsName
//	@param webFsName
//	@param httpServer
//	@return error
func InitHServer(httpServer *restapp.HttpServer, fileName string, srcFsName string, webFsName string, env *env.Env, autoRestart bool, options ...InitOptions) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
		if err != nil {
			err = errors.NewErr(err, "hserver.InitHServer()")
		}
	}()
	fact := NewFactory()
	hServer := env.App.HServer
	if !hServer.Enable {
		return nil
	}

	srcFs, ok := env.Fsm.GetFs(srcFsName)
	if !ok {
		return fmt.Errorf("%s fs not exists", srcFsName)
	}

	var webFs afero.Fs
	if webFsName != "" {
		webFs, ok = env.Fsm.GetFs(webFsName)
		if !ok {
			return fmt.Errorf(" %s fs not exists", webFsName)
		}
	}

	server, err := fact.NewServer(httpServer, fileName, srcFs, fact, env)
	if err != nil {
		return err
	}
	for _, option := range options {
		if err := option(server); err != nil {
			return err
		}
	}
	vData := map[string]any{
		"server": server,
		"pkg":    server.Pkg().Items(),
	}
	irisApp := httpServer.App()

	if webFs != nil {
		nodeModulesFs := env.Fsm.GetFsByTag("node-modules")
		cfg := &file_handler.Config{
			ServerFs:    srcFs,
			WebFs:       webFs,
			NodeModules: nodeModulesFs,
			Env:         env,
		}
		fileHandler := file_handler.NewHandler(irisApp, vData, cfg)
		httpServer.App().Get("/{file:path}", fileHandler.Handle)
		restapi.View = fileHandler.RenderView
	}
	if srcFs != nil && autoRestart {
		NewWatcher(server, srcFs, func(rootPath, fileName string, eventType fs.WatcherEventType) error {
			server.Logs(logrus.InfoLevel, "server.restart()")
			return restapp.Restart()
		})
	}

	return server.Start()
}
