package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/drawio/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/irisutils"
)

type DrawIoAPI struct {
	env     *env.Env
	service *service.FileService
}

type DrawIoAPISaveFileRequest struct {
	XML  string         `json:"xml"`
	Diff map[string]any `json:"diff"`
}

func NewDrawIoAPI(env *env.Env) *DrawIoAPI {
	service := service.NewFileService()
	return &DrawIoAPI{
		env:     env,
		service: service,
	}
}

func (s *DrawIoAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/drawio/file/save", "SaveFile")
	b.Handle("GET", "/drawio/file/read", "ReadFile")
}

func (s *DrawIoAPI) ReadFile(ctx iris.Context) {
	gp.Try(func() error {
		fileName := ctx.URLParamDefault("file", "")
		content, err := s.service.Read(fileName)
		if err == nil {
			_, err = ctx.Write(content)
			ctx.StatusCode(iris.StatusOK)
		}
		return err
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}

func (s *DrawIoAPI) SaveFile(ctx iris.Context) {
	gp.Try(func() error {
		var request DrawIoAPISaveFileRequest
		err := ctx.ReadJSON(&request)
		if err != nil {
			return err
		}

		fileName := ctx.URLParamDefault("file", "")
		return s.service.Save(fileName, []byte(request.XML))
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}
