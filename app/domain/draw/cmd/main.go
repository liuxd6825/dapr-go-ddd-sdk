package main

import (
	graphapi "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	masterapi "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp/cmd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
)

var (
	Version   = "1.0.0"
	BuildTime = ""
	GitHead   = ""
)

func main() {
	appcmd.StartApp(&appcmd.AppStartOptions{
		AppTitle:  "drawio服务器",
		Version:   Version,
		BuildTime: BuildTime,
		GitHead:   GitHead,
		Actors:    nil,
		OnHServerInitEvent: func(server element.Server) error {
			drawio := types.NewCMap[any]()
			drawio.Add("fileService", service.NewFileService())
			server.Pkg().Add("drawio", drawio)
			return nil
		},
		OnInitEvent: func(server *restapp.HttpServer) error {
			return nil
		},
		OnStartEvent: func(server *restapp.HttpServer) error {
			masterapi.RegisterAllApi(server.App(), "/api/v1", server.EnvConfig(), "")
			graphapi.RegisterAllApi(server.App(), "/api/v1", server.EnvConfig(), "")
			return nil
		},
	})
}
