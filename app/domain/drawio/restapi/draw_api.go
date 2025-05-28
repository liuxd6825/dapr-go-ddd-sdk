package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/restapi/request"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/drawio/service/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
	"net/url"
)

type DrawAPI struct {
	env          *env.Env
	fileService  *service.FileService
	graphService *service.GraphService
	drawDao      *dao.DrawDao
}

func NewDrawIoAPI(env *env.Env, rootPath string) *DrawAPI {
	fileService := service.NewFileService()
	graphService := service.NewGraphService()
	return &DrawAPI{
		env:          env,
		fileService:  fileService,
		graphService: graphService,
		drawDao:      dao.NewDrawDao(),
	}
}

func (s *DrawAPI) BeforeActivation(b mvc.BeforeActivation) {

	b.Handle(iris.MethodPost, "/draw", "Create")
	b.Handle(iris.MethodPut, "/draw", "Update")
	b.Handle(iris.MethodDelete, "/draw:deleteBatch", "DeleteByIds")
	b.Handle(iris.MethodDelete, "/draw", "DeleteById")

	b.Handle(iris.MethodGet, "/draw/{id}", "FindById")
	b.Handle(iris.MethodGet, "/draw", "FindPaging")

	b.Handle(iris.MethodPost, "/draw:file/{id}", "SaveFile")
	b.Handle(iris.MethodGet, "/draw:file/{id}", "ReadFile")

}

func (s *DrawAPI) Create(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd command.CreateCommand
		if err := web.GetCommandPost(ictx, &cmd); err != nil {
			return err
		}
		draw := &model.Draw{}
		draw.Id = cmd.Data.Id
		draw.CaseId = cmd.Data.CaseId
		draw.TenantId = appctx.GetTenantId2(ctx)
		draw.Name = cmd.Data.Name
		draw.FileName = cmd.Data.Id
		draw.Status = cmd.Data.Status
		draw.Remark = cmd.Data.Remark
		res := s.drawDao.Create(ctx, draw)
		if res.RowsAffected > 0 {
			_ = ictx.JSON(draw)
		}
		return s.fileService.Create(ctx, draw.CaseId, draw.Id)
	})
}

func (s *DrawAPI) Update(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		var cmd command.UpdateCommand
		if err := web.GetCommandPost(ictx, &cmd); err != nil {
			return err
		}
		draw := &model.Draw{}
		draw.Id = cmd.Data.Id
		draw.CaseId = cmd.Data.CaseId
		draw.TenantId = appctx.GetTenantId2(ctx)
		draw.Name = cmd.Data.Name
		draw.FileName = cmd.Data.Id
		draw.Status = cmd.Data.Status
		draw.Remark = cmd.Data.Remark
		res := s.drawDao.Update(ctx, draw)
		if res.RowsAffected > 0 {
			return ictx.JSON(draw)
		}
		return nil
	})
}

func (s *DrawAPI) DeleteById(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().GetString("id")
		res := s.drawDao.DeleteById(ctx, id)
		if res.RowsAffected == 0 {
			return errors.ErrorOf("delete 0 by id: %s", id)
		}
		return nil
	})
}

func (s *DrawAPI) DeleteByIds(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().GetString("id")
		res := s.drawDao.DeleteById(ctx, id)
		if res.RowsAffected == 0 {
			return errors.ErrorOf("delete 0 by id: %s", id)
		}
		return nil
	})
}

func (s *DrawAPI) FindPaging(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().GetString("id")
		query, err := web.GetFindPagingRequest(ictx)
		if err != nil {
			return err
		}
		draw := s.drawDao.FindPaging(ctx, query)
		if draw != nil {
			return web.SetData(ictx, draw)
		}
		return errors.ErrorOf("delete 0 by id: %s", id)
	})
}

func (s *DrawAPI) FindById(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		id := ictx.Params().GetString("id")
		draw := s.drawDao.FindById(ctx, id)
		if draw != nil {
			return web.SetData(ictx, draw)
		}
		return errors.ErrorOf("delete 0 by id: %s", id)
	})
}

func (s *DrawAPI) ReadFile(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		draw := s.getDrawById(ctx, ictx.Params().GetString("id"))
		ictx.Header("Case_Id", draw.CaseId)
		ictx.Header("Title", url.QueryEscape(draw.Name))
		content, err := s.fileService.Read(ctx, draw.CaseId, draw.FileName)
		if err == nil {
			_, err = ictx.Write(content)
		}
		return err
	})
}

func (s *DrawAPI) SaveFile(ictx iris.Context) {
	web.Try(ictx, func(ctx context.Context) error {
		draw := s.getDrawById(ctx, ictx.Params().GetString("id"))
		var saveRequest request.SaveFileRequest
		err := ictx.ReadJSON(&saveRequest)
		if err != nil {
			return err
		}
		err = s.fileService.Save(ctx, draw.CaseId, draw.FileName, saveRequest.XML)
		if err != nil {
			return err
		}
		return s.graphService.Save(ctx, &saveRequest)
	})
}

func (s *DrawAPI) getDrawById(ctx context.Context, id string) *model.Draw {
	if id == "" {
		panic(errors.New("id is required"))
	}

	draw := s.drawDao.FindById(ctx, id)
	if draw == nil {
		panic(errors.New("can not find draw"))
	}
	return draw
}
