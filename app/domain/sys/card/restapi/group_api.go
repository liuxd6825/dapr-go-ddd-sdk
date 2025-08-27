package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type GroupAPI struct {
	env          *env.Env
	groupService *service.GroupService
	rootPath     string
}

func NewGroupAPI(env *env.Env, rootPath string) *GroupAPI {
	return &GroupAPI{
		env:          env,
		groupService: service.NewGroupService(),
		rootPath:     rootPath,
	}
}

func (s *GroupAPI) InitController(app *iris.Application) error {
	s.groupService = service.NewGroupService()
	ctl := restapi.NewController(app, s.rootPath+"/card", s)
	ctl.Post("/group", "Create")
	ctl.Put("/group", "Update")
	ctl.Delete("/group", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/group/{id}", "FindById")
	ctl.GetPaging("/group", "FindPaging")
	return nil
}

func (s *GroupAPI) Create(ctx context.Context, cmd *command.GroupCreateCommand) error {
	return s.groupService.Create(ctx, cmd)
}

func (s *GroupAPI) Update(ctx context.Context, cmd *command.GroupUpdateCommand) error {
	return s.groupService.Update(ctx, cmd)
}

func (s *GroupAPI) Delete(ctx context.Context, cmd *command.GroupDeleteCommand) error {
	return s.groupService.Delete(ctx, cmd)
}

func (s *GroupAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Group, error) {
	return s.groupService.FindById(ctx, qry)
}

func (s *GroupAPI) FindPaging(ctx context.Context, qry *query.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Group], error) {
	return s.groupService.FindPaging(ctx, qry.CaseId, qry)
}
