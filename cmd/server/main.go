package main

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/jsserver"
	"os"
)

func main() {
	rootPath, _ := os.Getwd()
	app := iris.New()
	server := jsserver.New(app, rootPath+"/src", true)
	if err := server.Run(); err != nil {
		panic(err)
	}
	if err := app.Run(iris.Addr(":8080")); err != nil {
		panic(err)
	}
}
