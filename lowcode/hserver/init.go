package hserver

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/ctx_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/feign_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/json_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/logs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/params_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/schema_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/strings_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/tpl_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
)

// InitServer
//
//	@Description: 添加到restapp的初始化函数 Options.Init
//	@param s
//	@return error
func InitServer(fileName string, srcFsName string, tplFsName string, httpServer *restapp.HttpServer) error {
	env := httpServer.EnvConfig()
	if !env.App.RsServer.Enable {
		return nil
	}
	envCfg := httpServer.EnvConfig()

	srcFs, err := httpServer.EnvConfig().GetFs(srcFsName)
	if err != nil {
		return fmt.Errorf("srcFs %s not exists", srcFsName)
	}

	tplFs, err := httpServer.EnvConfig().GetFs(tplFsName)
	if err != nil {
		return fmt.Errorf("tplFs %s not exists", tplFsName)
	}

	server, err := NewServer(httpServer.App(), fileName, srcFs, envCfg)
	pkg := types.NewCMap[any]()

	pkg.Set("mongo", mongodb.New(envCfg))
	pkg.Set("template", tpl_pkg.New(envCfg, server, tplFs))
	pkg.Set("feign", feign_pkg.New(server))
	pkg.Set("fs", server.fsPkg)
	pkg.Set("context", ctx_pkg.New())
	pkg.Set("schema", schema_pkg.New(server))
	pkg.Set("params", params_pkg.New(server))
	pkg.Set("env", envCfg)
	pkg.Set("logs", logs_pkg.New())
	pkg.Set("json", json_pkg.New())
	pkg.Set("strings", strings_pkg.New())

	server.SetPkg(pkg)

	server.SetRunValues(map[string]any{
		"console": NewConsole(logrus.New()),
	})

	if err != nil {
		return err
	}

	return server.Start()

}
