package restapi

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
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

	b.Handle("GET", "/schema:properties/{file:path}", "GetSchemaProperties")
	b.Handle("GET", "/schema:pathAll/{path:path}", "ReadAllPath")
	b.Handle("POST", "/schema:createPath", "CreatePath")
	b.Handle("POST", "/schema:createFile", "CreateFile")
	b.Handle("PUT", "/schema:renameFile", "RenameFile")
	b.Handle("DELETE", "/schema:removeFile", "RemoveFile")
	b.Handle("DELETE", "/schema:removePath", "RemovePath")
	b.Handle("PUT", "/schema:save/{file:path}", "WriteJson")
}

func (s *SchemaAPI) ReadFile(ctx iris.Context, file string) {
	gp.Try(func() error {
		data, err := s.service.ReadFile(ctx, file)
		if err == nil {
			_, err = ctx.Write(data)
		}
		return err
	}).Catch(func(err error) {
		web.SetError(ctx, err)
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
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) WriteJson(ctx iris.Context, file string) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		return s.service.WriteJson(ctx, file, body)
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) ReadPath(ctx iris.Context, path string) {
	gp.Try(func() error {
		list := s.service.ReadPath(ctx, path)
		return ctx.JSON(list)
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) GetSchema(ctx iris.Context, file string) {
	gp.Try(func() error {
		sch := s.service.GetSchema(ctx, file)
		view := sch.GetView()
		return ctx.JSON(view)
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) GetSchemaProperties(ctx iris.Context, file string) {
	gp.Try(func() error {
		sch := s.service.GetSchema(ctx, file)
		view := sch.GetView()
		return ctx.JSON(view.Properties)
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) ReadAllPath(ctx iris.Context, path string) {
	gp.Try(func() error {
		list := s.service.ReadAllPath(ctx, path)
		return ctx.JSON(list)
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) CreatePath(ctx iris.Context) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		data := map[string]string{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		pathName := data["pathName"]
		err = s.service.CreatePath(ctx, pathName)
		return err
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) CreateFile(ctx iris.Context) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		data := map[string]string{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		fileName := data["fileName"]
		err = s.service.CreateFile(ctx, fileName)
		return err
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) RenameFile(ctx iris.Context) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		data := map[string]string{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		oldFileName := data["oldFileName"]
		newFileName := data["newFileName"]
		err = s.service.RenameFile(ctx, oldFileName, newFileName)
		return err
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) RemoveFile(ctx iris.Context) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		data := map[string]string{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		fileName := data["fileName"]
		err = s.service.RemoveFile(ctx, fileName)
		return err
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}

func (s *SchemaAPI) RemovePath(ctx iris.Context) {
	gp.Try(func() error {
		body, err := ctx.GetBody()
		if err != nil {
			return err
		}
		data := map[string]string{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return err
		}
		pathName := data["pathName"]
		err = s.service.RemovePath(ctx, pathName)
		return err
	}).Catch(func(err error) {
		web.SetError(ctx, err)
	})
}
