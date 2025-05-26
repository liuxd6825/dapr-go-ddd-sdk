package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/drawio/restapi/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/drawio/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/irisutils"
)

type DrawIoAPI struct {
	env          *env.Env
	fileService  *service.FileService
	graphService *service.GraphService
}

func NewDrawIoAPI(env *env.Env) *DrawIoAPI {
	fileService := service.NewFileService()
	graphService := service.NewGraphService()
	return &DrawIoAPI{
		env:          env,
		fileService:  fileService,
		graphService: graphService,
	}
}

func (s *DrawIoAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/drawio/file/save", "SaveFile")
	b.Handle("GET", "/drawio/file/read", "ReadFile")
}

func (s *DrawIoAPI) ReadFile(ctx iris.Context) {
	gp.Try(func() error {
		fileName := ctx.URLParamDefault("file", "")
		content, err := s.fileService.Read(fileName)
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
		var saveRequest request.SaveFileRequest
		err := ctx.ReadJSON(&saveRequest)
		if err != nil {
			return err
		}

		fileName := ctx.URLParamDefault("file", "")
		err = s.fileService.Save(fileName, saveRequest.XML)
		if err != nil {
			return err
		}

		return s.graphService.Save(ctx, &saveRequest)
	}).Catch(func(err error) {
		irisutils.SetError(ctx, err)
	})
}
