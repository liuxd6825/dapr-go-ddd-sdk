package restapi

import (
	"context"
	"fmt"
	"net/url"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service/model"
	service2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type DrawFileRequest struct {
	Id string `json:"id" path:"id" required:"true" title:"案件ID"`
}

type DrawAPI struct {
	env          *env.Env
	fileService  *service.FileService
	graphService *service2.GraphService
	drawDao      *dao.DrawDao
	rootPath     string
}

func NewDrawIoAPI(env *env.Env, rootPath string) *DrawAPI {
	fileService := service.NewFileService()
	graphService := service2.NewGraphService()
	return &DrawAPI{
		env:          env,
		fileService:  fileService,
		graphService: graphService,
		drawDao:      dao.NewDrawDao(service.DBKey),
		rootPath:     rootPath,
	}
}

func (s *DrawAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/draw", "Create")
	b.Handle(iris.MethodPut, "/draw", "Update")
	b.Handle(iris.MethodDelete, "/draw:deleteBatch", "DeleteByIds")
	b.Handle(iris.MethodDelete, "/draw", "DeleteById")

	b.Handle(iris.MethodGet, "/draw/{id}", "FindById")
	b.Handle(iris.MethodGet, "/draw", "FindPaging")

	b.Handle(iris.MethodPost, "/draw/{id}/file", "SaveFile")
	b.Handle(iris.MethodGet, "/draw/{id}/file", "ReadFile")
}

func (s *DrawAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "draw.DrawAPI", s)
	controller.Post("/draw", "Create")
	controller.Put("/draw", "Update")
	controller.Delete("/draw:deleteBatch", "DeleteByIds", restapi.WithParamsInBody(true))
	controller.Delete("/draw", "DeleteById", restapi.WithParamsInBody(true))
	controller.GetOne("/draw/{id}", "FindById")
	controller.GetPaging("/draw", "FindPaging")
	controller.Post("/draw/{id}/file", "SaveFile")
	controller.GetData("/draw/{id}/file", "ReadFile")
	return controller
}

func (s *DrawAPI) Create(ctx context.Context, ictx iris.Context, cmd *command.CreateCommand) error {
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
	return s.fileService.Create(ctx, draw.CaseId, draw.Id, draw.FileName)
}

func (s *DrawAPI) Update(ctx context.Context, cmd *command.UpdateCommand) (any, error) {
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
		return res, nil
	}
	return nil, nil
}

func (s *DrawAPI) DeleteById(ctx context.Context, cmd *command.DeleteCommand) error {
	return s.drawDao.DeleteById(ctx, cmd.Data.Id).Error
}

func (s *DrawAPI) DeleteByIds(ctx context.Context, cmd *command.DeleteByIdsCommand) error {
	return s.drawDao.DeleteByIds(ctx, cmd.Data.Ids).Error
}

func (s *DrawAPI) FindPaging(ctx context.Context, query *restapi.FindPagingByCaseIdRequest) any {
	query.SetMustFilter(fmt.Sprintf("case_id=='%s'", query.GetCaseId()))
	return s.drawDao.FindPaging(ctx, query)
}

func (s *DrawAPI) FindById(ctx context.Context, query *restapi.FindByIdRequest) (any, error) {
	draw, err := s.drawDao.FindById(ctx, query.Id)
	return draw, err
}

func (s *DrawAPI) ReadFile(ctx context.Context, ictx iris.Context, params *DrawFileRequest) error {
	draw, err := s.getDrawById(ctx, params.Id)
	if err != nil {
		return err
	}
	ictx.Header("Case_Id", draw.CaseId)
	ictx.Header("Title", url.QueryEscape(draw.Name))
	content, err := s.fileService.Read(ctx, draw.CaseId, draw.FileName)
	if err == nil {
		_, err = ictx.Write(content)
	}
	return err
}

func (s *DrawAPI) SaveFile(ctx context.Context, cmd *command.DrawSaveFileCommand) error {

	id := cmd.Data.DrawId
	draw, err := s.getDrawById(ctx, id)
	if err != nil {
		return err
	}
	if draw == nil {
		return errors.ErrorOf("没有找到分析图: %s", id)
	}

	err = s.fileService.Save(ctx, cmd.GetCaseId(), cmd.GetDrawId(), cmd.GetFileName(), cmd.GetXML())
	if err != nil {
		return err
	}
	return s.graphService.Save(ctx, cmd)

}

func (s *DrawAPI) getDrawById(ctx context.Context, id string) (*model.Draw, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	draw, err := s.drawDao.FindById(ctx, id)
	return draw, err
}
