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
	cardService  *service.CardService
	rootPath     string
}

func NewGroupAPI(env *env.Env, rootPath string) *GroupAPI {
	return &GroupAPI{
		env:          env,
		groupService: service.NewGroupService(),
		cardService:  service.NewCardService(),
		rootPath:     rootPath,
	}
}

func (s *GroupAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	s.groupService = service.NewGroupService()
	ctl := restapi.NewController(app, s.rootPath+"/card", "sys.GroupAPI", s)
	ctl.Post("/group", "Create")
	ctl.Put("/group", "Update")
	ctl.Delete("/group", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/group/{id}", "FindById")
	ctl.GetPaging("/group", "FindPaging")
	ctl.GetData("/group:view", "FindGroupViewById")
	ctl.GetData("/group:home-id", "FindByHomeId")
	return ctl
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

func (s *GroupAPI) FindByHomeId(ctx context.Context, qry *query.FindByHomeIdQuery) ([]*model.Group, error) {
	return s.groupService.FindByHomeId(ctx, qry.HomeId)
}

func (s *GroupAPI) FindGroupViewById(ctx context.Context, qry *query.FindByIdQuery) (*model.GroupView, error) {
	g, err := s.groupService.FindById(ctx, qry)
	if err != nil {
		return nil, err
	}

	cards, err := s.cardService.FindByGroupId(ctx, qry.Id)
	if err != nil {
		return nil, err
	}

	groupView := &model.GroupView{
		Group: g,
		Cards: cards,
	}

	return groupView, nil
}
