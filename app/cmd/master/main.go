package main

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis"
	doc "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/document/restapi"
	drawio "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph"
	excelImport "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/metrics"
	rag "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/restapi"
	card "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/restapi"
	tag "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/tag/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp/cmd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/tasks"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
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
			drawioPkg := types.NewCMap[any]()
			drawioPkg.Add("fileService", service.NewFileService())
			server.Pkg().Add("drawio", drawioPkg)
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
				card.RegisterAllApi(app, baseUrl, env)
				metrics.RegisterAllApi(app, baseUrl, env)
				tasks.RunWorker()
				return nil
			})
			return err
		},
	})
}
