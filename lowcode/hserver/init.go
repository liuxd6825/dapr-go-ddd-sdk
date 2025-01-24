package hserver

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/handler/file_handler"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/console_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/ctx_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/feign_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/html_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/json_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/logs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/params_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/schema_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/strings_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/tpl_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
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
func InitServer(fileName string, srcFsName string, webFsName string, httpServer *restapp.HttpServer) error {
	fact := NewFactory()
	env := httpServer.EnvConfig()
	if !env.App.RsServer.Enable {
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
	server, err := fact.NewServer(httpServer.App(), fileName, srcFs, fact, envCfg)
	if err != nil {
		return err
	}
	pkg := types.NewCMap[any]()

	pkg.Set("mongo", mongodb.New(envCfg))
	if webFs != nil {
		pkg.Set("template", tpl_pkg.New(envCfg, server, webFs))
	}
	pkg.Set("feign", feign_pkg.New(server))
	pkg.Set("fs", server.FsPkg())
	pkg.Set("context", ctx_pkg.New())
	pkg.Set("schema", schema_pkg.New(server))
	pkg.Set("params", params_pkg.New(server))
	pkg.Set("env", envCfg)
	pkg.Set("logs", logs_pkg.New())
	pkg.Set("json", json_pkg.New())
	pkg.Set("strings", strings_pkg.New())
	pkg.Set("html", html_pkg.New(server))

	server.SetPkg(pkg)

	server.SetRunValues(map[string]any{
		"console": console_pkg.NewConsole(logrus.New()),
	})

	vData := map[string]any{
		"server": server,
		"pkg":    server.Pkg().Items(),
	}
	vApp := httpServer.App()

	if webFs != nil {
		fileHandler := file_handler.NewHandler(webFs, vApp, vData)
		httpServer.App().Get("/{file:path}", fileHandler.Handle)
		//NewWatcher(server, webFs)
	}
	if srcFs != nil {
		NewWatcher(server, srcFs, func(rootPath, fileName string, eventType fs.WatcherEventType) error {
			server.Logs(logrus.InfoLevel, "server.restart()")
			return server.Restart()
		})
	}

	return server.Start()
}
