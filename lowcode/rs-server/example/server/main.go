package main

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server"
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/localfs"
	"github.com/spf13/afero"
	"os"
	"strings"
)

func main() {
	app := iris.New()

	fs, err := newServerFsConfig()
	if err != nil {
		panic(err)
	}

	server, err := rs_server.NewServer(app, nil, fs, nil, "./main.js", true)
	if err != nil {
		panic(err)
	}
	if err = server.Run(); err != nil {
		panic(err)
	}
	if err = app.Run(iris.Addr(":8080")); err != nil {
		panic(err)
	}
}

func newServerFsConfig() (*rs_server.FsConfig, error) {
	// 获取当前目录 必须是项目下的lowcode/rs-server/example/server目录
	rootPath, _ := os.Getwd()
	if !strings.HasSuffix(rootPath, "/example/server") {
		rootPath = rootPath + "/lowcode/rs-server/example/server"
	}
	// 设置源代码目录为当前目录下的src目录
	fileFs, err := localfs.NewFs(localfs.Config{Id: "file", Path: rootPath + "/src"})
	if err != nil {
		return nil, err
	}
	fsc := &rs_server.FsConfig{
		FileFs: fileFs,
		HttpFs: afero.NewMemMapFs(),
	}
	return fsc, nil
}

/*
func newRenderFsManger(tplPath string) *fs.Manager {

}

func htmlHandler(app *iris.Application, tplPath string) {
	if render, err := template.NewHandler(newRenderFsManger(tplPath)); err != nil {
		panic(err)
	} else {
		app.Get("/api/v1/html/{reposName}/{filePath:path}", func(ictx iris.Context) {
			reposName := ictx.URLParam("reposName")
			filePath := ictx.URLParam("filePath")
			render(ictx, reposName, filePath)
		})
	}
}
*/
