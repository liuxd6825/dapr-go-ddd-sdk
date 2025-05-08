package file_handler

import (
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
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
	handler := NewHandler(fs, app)
	app.Handle("GET", "/{file:path}", handler.Handle)
	err = app.Run(iris.Addr(":8080"))
	if err != nil {
		t.Fatal(err)
	}

}
