package main

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/master/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp/cmd"
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
		OnInitEvent: func(server *restapp.HttpServer) error {
			return nil
		},
		OnStartEvent: func(server *restapp.HttpServer) error {
			restapi.RegisterDrawIo(server.App(), "/api/v1", server.EnvConfig(), "")
			return nil
		},
	})
}
