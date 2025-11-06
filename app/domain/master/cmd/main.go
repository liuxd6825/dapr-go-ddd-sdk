package main

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp"
	appcmd "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/core/restapp/cmd"
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
		OnInitEvent: func(server *restapp.HttpServer) error {
			return nil
		},
		OnStartEvent: func(server *restapp.HttpServer) error {
			//restapi.RegisterAllApi(server.App(), "/api/v1", server.EnvConfig(), "")
			return nil
		},
	})
}
