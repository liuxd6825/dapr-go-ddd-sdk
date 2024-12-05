package hserver

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/ctx_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/feign_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/params_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/schema_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/tpl_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
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

	server.SetPkg(map[string]any{
		"mongo":    mongodb.New(envCfg),
		"template": tpl_pkg.New(envCfg, server, tplFs),
		"feign":    feign_pkg.New(server),
		"fs":       server.fsm,
		"context":  ctx_pkg.New(),
		"schema":   schema_pkg.New(server),
		"params":   params_pkg.New(server),
		"env":      envCfg,
	})

	server.SetRunValues(map[string]any{
		"console": NewConsole(logrus.New()),
	})

	if err != nil {
		return err
	}

	return server.Start()

}
