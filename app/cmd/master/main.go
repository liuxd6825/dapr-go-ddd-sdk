package main

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis"
	doc "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/restapi"
	drawio "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/restapi"
	draw_service "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph"
	excelImport "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/metrics"
	notify "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/notify/restapi"
	rag "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys"
	sys_code_service "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/code/service"
	tag "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp/cmd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/tasks"
)

var (
	Version   = "1.0.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	appcmd.StartApp(&appcmd.AppStartOptions{
		AppTitle:  "master服务器",
		Version:   Version,
		BuildTime: BuildTime,
		GitHead:   GitHead,
		Actors:    nil,
		OnHServerInitEvent: func(server element.Server) error {
			drawPkg := map[string]any{}
			drawPkg["fileService"] = draw_service.NewFileService()
			server.Pkg().Add("drawio", drawPkg)

			sysPkg := map[string]any{}
			sysPkg["codeService"] = sys_code_service.NewCodeService()
			server.Pkg().Add("sys", sysPkg)

			return nil
		},
		OnInitEvent: func(server *restapp.HttpServer) error {
			return nil
		},
		OnStartEvent: func(server *restapp.HttpServer) error {
			err := logs.DebugStart(context.Background(), logs.Fields{"service name ": "all"}, func() error {
				baseUrl := "/api/v1.0"
				app := server.App()
				env := server.EnvConfig()
				master.Init(app, baseUrl, env)
				analysis.Init(app, baseUrl, env)
				drawio.RegisterAllApi(app, baseUrl, env)
				graph.RegisterAllApi(app, baseUrl, env)
				rag.RegisterAllApi(app, baseUrl, env)
				doc.RegisterAllApi(app, baseUrl, env)
				tag.RegisterAllApi(app, baseUrl, env)
				excelImport.RegisterAllApi(app, baseUrl, env)
				sys.RegisterAllApi(app, baseUrl, env)
				metrics.RegisterAllApi(app, baseUrl, env)
				notify.RegisterAllApi(app, baseUrl, env)
				/*
					if err := lowcode.Run(server); err != nil {
						return err
					}
				*/
				tasks.RunWorker()
				return nil
			})
			return err
		},
	})
}
