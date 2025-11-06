package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/gp"
)

type HtmlAPI struct {
	service service.IFileService
}

func NewHtmlAPI(env *env.Env, rootPath string) *HtmlAPI {
	return &HtmlAPI{
		service: service.NewFileService(env, rootPath),
	}
}

func (s *HtmlAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/html:path/{path:path}", "ReadPath")
	b.Handle("GET", "/html:pathAll/{path:path}", "ReadAllPath")
	b.Handle("POST", "/html:save", "WriteFile")
}

func (s *HtmlAPI) ReadFile(ctx iris.Context, file string) {
	gp.Try(func() error {
		data, err := s.service.ReadFile(ctx, file)
		if err == nil {
			_, err = ctx.Write(data)
		}
		return err
	}).Catch(func(err error) {
		restapi.SetError(ctx, err)
	})
}

func (s *HtmlAPI) WriteFile(ctx iris.Context) {
	gp.Try(func() error {
		fileName := ctx.FormValue("fileName")
		data := ctx.FormValue("data")
		return s.service.WriteFile(ctx, fileName, []byte(data))
	}).Catch(func(err error) {
		restapi.SetError(ctx, err)
	})
}

func (s *HtmlAPI) ReadPath(ctx iris.Context, path string) {
	gp.Try(func() error {
		list := s.service.ReadPath(ctx, path, false)
		return ctx.JSON(list)
	}).Catch(func(err error) {
		restapi.SetError(ctx, err)
	})
}

func (s *HtmlAPI) ReadAllPath(ctx iris.Context, path string) {
	gp.Try(func() error {
		list := s.service.ReadPath(ctx, path, true)
		return ctx.JSON(list)
	}).Catch(func(err error) {
		restapi.SetError(ctx, err)
	})
}
