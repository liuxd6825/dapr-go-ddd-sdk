package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/restapi/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
)

type DrawIoAPI struct {
	env          *env.Env
	fileService  *service.FileService
	graphService *service.GraphService
}

func NewDrawIoAPI(env *env.Env, rootPath string) *DrawIoAPI {
	fileService := service.NewFileService()
	graphService := service.NewGraphService()
	return &DrawIoAPI{
		env:          env,
		fileService:  fileService,
		graphService: graphService,
	}
}

func (s *DrawIoAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("POST", "/drawio-file:save", "SaveFile")
	b.Handle("GET", "/drawio-file:read", "ReadFile")
}

func (s *DrawIoAPI) ReadFile(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		fileName := ictx.URLParamDefault("file", "")
		content, err := s.fileService.Read(fileName)
		if err == nil {
			_, err = ictx.Write(content)
		}
		return err
	})
}

func (s *DrawIoAPI) SaveFile(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var saveRequest request.SaveFileRequest
		err := ictx.ReadJSON(&saveRequest)
		if err != nil {
			return err
		}

		fileName := ictx.URLParamDefault("file", "")
		err = s.fileService.Save(fileName, saveRequest.XML)
		if err != nil {
			return err
		}

		return s.graphService.Save(ctx, &saveRequest)
	})
}
