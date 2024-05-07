package server

import (
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader"
)

func AddHtmlHandler(app *iris.Application, relativePath string) error {
	app.Get(relativePath+"/{repoName}/{filepath:path}", htmlHandler)
	return nil
}

// htmlHandler
//
//	@Description:
//	@receiver a
//	@param ictx
func htmlHandler(ictx iris.Context) {
	repoName := ictx.URLParam("repoName")
	filepath := ictx.URLParam("filepath")
	Do(ictx, func(ctx iris.Context) error {
		fileService, ok := loader.GetLoader(repoName)
		if !ok {
			return errors.New(fmt.Sprintf("respName %s not ", repoName))
		}
		return template.Render(ictx, ictx.ResponseWriter(), fileService, filepath)
	})
}
