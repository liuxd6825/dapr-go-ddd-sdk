package hserver

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/dapr"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/ctx_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/feign_pkg"
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
func InitServer(httpServer *restapp.HttpServer) error {
	env := httpServer.EnvConfig()
	if !env.App.RsServer.Enable {
		return nil
	}
	envCfg := httpServer.EnvConfig()
	server, err := NewServer(httpServer.App(), "/src-server/xsrc/server.html", envCfg, func(server *Server) {
		data := map[string]any{
			"DAPR_HOST":      envCfg.GetDaprHost,
			"DAPR_HTTP_PORT": envCfg.GetDaprHttpPort,
			"DAPR_GRPC_PORT": envCfg.GetDaprGrpcPort,
			"APP_ID":         envCfg.GetAppId,
			"APP_NAME":       envCfg.GetAppName,
			"env":            envCfg,
			"daprClient":     dapr.GetDaprClient(),
			"console":        newConsole(logrus.New()),
			"db":             db_pkg.New(envCfg),
			"template":       tpl_pkg.New(envCfg, server),
			"feign":          feign_pkg.New(server),
			"fs":             server.fs,
			"context":        ctx_pkg.New(),
			"schemas":        schema_pkg.New(server),
		}
		server.SetRunValues(data)
	})

	if err != nil {
		return err
	}

	return server.Run()

}
