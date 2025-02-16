package hserver

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/handler/file_handler"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// InitServer
//
//	@Description:
//	@param fileName
//	@param srcFsName
//	@param webFsName
//	@param httpServer
//	@return error
func InitServer(fileName string, srcFsName string, webFsName string, httpServer *restapp.HttpServer, autoRestart bool) error {
	fact := NewFactory()
	env := httpServer.EnvConfig()
	if !env.App.HServer.Enable {
		return nil
	}
	envCfg := httpServer.EnvConfig()

	srcFs, err := httpServer.EnvConfig().GetFs(srcFsName)
	if err != nil {
		return fmt.Errorf("%s fs not exists", srcFsName)
	}

	var webFs afero.Fs
	if webFsName != "" {
		webFs, err = httpServer.EnvConfig().GetFs(webFsName)
		if err != nil {
			return fmt.Errorf(" %s fs not exists", webFsName)
		}
	}
	fmt.Println("lowcode HServer")
	server, err := fact.NewServer(httpServer, fileName, srcFs, fact, envCfg)
	if err != nil {
		return err
	}
	vData := map[string]any{
		"server": server,
		"pkg":    server.Pkg().Items(),
	}
	vApp := httpServer.App()

	if webFs != nil {
		fileHandler := file_handler.NewHandler(webFs, vApp, vData)
		httpServer.App().Get("/{file:path}", fileHandler.Handle)
	}
	if srcFs != nil && autoRestart {
		NewWatcher(server, srcFs, func(rootPath, fileName string, eventType fs.WatcherEventType) error {
			server.Logs(logrus.InfoLevel, "server.restart()")
			return restapp.Restart()
		})
	}

	return server.Start()
}
