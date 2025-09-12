package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type TemplateAPI struct {
	env             *env.Env
	templateService *service.TemplateService
	rootPath        string
}

func NewTemplateAPI(env *env.Env, rootPath string) *TemplateAPI {
	return &TemplateAPI{
		env:             env,
		templateService: service.NewTemplateService(),
		rootPath:        rootPath,
	}
}

func (s *TemplateAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.templateService = service.NewTemplateService()
	ctl := restapi.NewController(app, s.rootPath+"/import", "TemplateAPI", s)
	ctl.Post("/template", "Create")
	ctl.Put("/template", "Update")
	ctl.Delete("/template", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/template/{id}", "FindById")
	ctl.GetPaging("/template", "FindPaging")
	return ctl
}

func (s *TemplateAPI) Create(ctx context.Context, cmd *command.TempCreateCommand) error {
	return s.templateService.Create(ctx, cmd)
}

func (s *TemplateAPI) Update(ctx context.Context, cmd *command.TempUpdateCommand) error {
	return s.templateService.Update(ctx, cmd)
}

func (s *TemplateAPI) Delete(ctx context.Context, cmd *command.TempDeleteCommand) error {
	return s.templateService.Delete(ctx, cmd)
}

func (s *TemplateAPI) FindById(ctx context.Context, qry *query.TemplateFindByIdQuery) (*model.Template, error) {
	return s.templateService.FindById(ctx, qry)
}

func (s *TemplateAPI) FindPaging(ctx context.Context, qry *query.TemplateFindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Template], error) {
	return s.templateService.FindPaging(ctx, qry.CaseId, qry)
}
