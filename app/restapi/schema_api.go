package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/irisutils"
)

type SchemaAPI struct {
	service service.ISchemaService
}

func NewSchemaAPI(env *env.Env, rootPath string) *SchemaAPI {
	return &SchemaAPI{
		service: service.NewSchemaService(env, rootPath),
	}
}

func (s *SchemaAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/schema:file/{file:path}", "ReadFile")
	b.Handle("GET", "/schema:path/{path:path}", "ReadPath")
	b.Handle("GET", "/schema/{file:path}", "GetSchema")
}

func (s *SchemaAPI) ReadFile(ctx iris.Context, file string) {
	gp.Try(func() error {
		data, err := s.service.ReadFile(ctx, file)
		if err == nil {
			_, err = ctx.Write(data)
		}
		return err
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}

func (s *SchemaAPI) WriteFile(ctx iris.Context, file string) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		return s.service.WriteFile(ctx, file, body)
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}

func (s *SchemaAPI) ReadPath(ctx iris.Context, path string) {
	gp.Try(func() error {
		list := s.service.ReadPath(ctx, path)
		return ctx.JSON(list)
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}

func (s *SchemaAPI) GetSchema(ctx iris.Context, file string) {
	gp.Try(func() error {
		sch := s.service.GetSchema(ctx, file)
		view := sch.GetView()
		return ctx.JSON(view)
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}
