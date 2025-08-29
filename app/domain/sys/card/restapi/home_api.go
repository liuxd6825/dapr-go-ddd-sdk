package restapi

import (
	"context"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type HomeAPI struct {
	env          *env.Env
	homeService  *service.HomeService
	groupService *service.GroupService
	cardService  *service.CardService
	rootPath     string
}

func NewHomeAPI(env *env.Env, rootPath string) *HomeAPI {
	return &HomeAPI{
		env:          env,
		homeService:  service.NewHomeService(),
		groupService: service.NewGroupService(),
		cardService:  service.NewCardService(),
		rootPath:     rootPath,
	}
}

func (s *HomeAPI) InitController(app *iris.Application) error {
	s.homeService = service.NewHomeService()
	ctl := restapi.NewController(app, s.rootPath+"/card", s)
	ctl.Post("/home", "Create")
	ctl.Post("/home:submit", "Submit")
	ctl.Put("/home", "Update")
	ctl.Delete("/home", "Delete", restapi.WithParamsInBody(true))
	ctl.GetOne("/home/{id}", "FindById")
	ctl.GetPaging("/home", "FindPaging")
	ctl.GetData("/home:user-home", "FindUserHomeByCode")
	return nil
}

func (s *HomeAPI) Create(ctx context.Context, cmd *command.HomeCreateCommand) error {
	return s.homeService.Create(ctx, cmd)
}

func (s *HomeAPI) Update(ctx context.Context, cmd *command.HomeUpdateCommand) error {
	return s.homeService.Update(ctx, cmd)
}

func (s *HomeAPI) Delete(ctx context.Context, cmd *command.HomeDeleteCommand) error {
	return s.homeService.Delete(ctx, cmd)
}

func (s *HomeAPI) Submit(ctx context.Context, cmd *command.HomeSubmitCommand) error {
	var err error
	if cmd.Data.Home == nil {
		//存在用户首页，直接更新即可
		if _, err = s.cardService.UpdateMany(ctx, cmd.Data.Card.U); err != nil {
			return err
		}
		if _, err = s.cardService.CreateMany(ctx, cmd.Data.Card.I); err != nil {
			return err
		}
		if _, err = s.cardService.DeleteMany(ctx, cmd.Data.Card.D); err != nil {
			return err
		}
		if _, err = s.cardService.DeleteByListGroupId(ctx, cmd.Data.Group.D); err != nil {
			return err
		}
		if _, err = s.groupService.UpdateMany(ctx, cmd.Data.Group.U); err != nil {
			return err
		}
		if _, err = s.groupService.CreateMany(ctx, cmd.Data.Group.I); err != nil {
			return err
		}
		if _, err = s.groupService.DeleteMany(ctx, cmd.Data.Group.D); err != nil {
			return err
		}
	} else {
		//不存在用户首页
		if err = s.homeService.CreateByModel(ctx, cmd.Data.Home); err != nil {
			return err
		}
		if _, err = s.cardService.CreateMany(ctx, cmd.Data.Card.I); err != nil {
			return err
		}
		if _, err = s.groupService.CreateMany(ctx, cmd.Data.Group.I); err != nil {
			return err
		}
	}
	return nil
}

func (s *HomeAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Home, error) {
	return s.homeService.FindById(ctx, qry)
}

func (s *HomeAPI) FindPaging(ctx context.Context, qry *query.FindPagingByCaseIdQuery) (idao.FindPagingResult[*model.Home], error) {
	return s.homeService.FindPaging(ctx, qry.CaseId, qry)
}

func (s *HomeAPI) FindUserHomeByCode(ctx context.Context, qry *query.FindByCodeQuery) (*model.HomeView, error) {
	home, err := s.homeService.FindByCode(ctx, qry)
	if err != nil {
		return nil, err
	}
	if home == nil {
		home, err = s.homeService.FindDefaultByCode(ctx, qry.Code)
		if err != nil {
			return nil, err
		}
	}
	if home == nil {
		return nil, errors.New("没有找到首页数据")
	}

	groups, err := s.groupService.FindByHomeId(ctx, home.Id)
	if err != nil {
		return nil, err
	}

	cards, err := s.cardService.FindByHomeId(ctx, home.Id)
	if err != nil {
		return nil, err
	}

	hv := &model.HomeView{
		Home:   home,
		Groups: groups,
		Cards:  cards,
	}
	return hv, nil
}
