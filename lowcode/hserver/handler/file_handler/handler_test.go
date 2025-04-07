package file_handler

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/localfs"
	"os"
	"testing"
)

func TestHandler_Handle(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	app := iris.New()
	fs, err := localfs.NewFs(localfs.Config{
		Name: "file",
		Path: path + "/test",
	})
	if err != nil {
		t.Fatal(err)
	}
	data := make(map[string]any)
	cfg := &Config{
		SrcFs: fs,
	}

	handler := NewHandler(app, data, cfg)
	app.Handle("GET", "/{file:path}", handler.Handle)
	err = app.Run(iris.Addr(":8080"))
	if err != nil {
		t.Fatal(err)
	}

}
