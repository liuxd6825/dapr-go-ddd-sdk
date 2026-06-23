package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/human/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// HumanSuspectAmountAPI 嫌疑人金额REST API 控制器
type HumanSuspectAmountAPI struct {
	service  *service.HumanSuspectAmountService
	rootPath string
}

func NewHumanSuspectAmountAPI(rootPath string) *HumanSuspectAmountAPI {
	return &HumanSuspectAmountAPI{
		rootPath: rootPath,
		service:  service.NewHumanSuspectAmountService(),
	}
}

func (s *HumanSuspectAmountAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanSuspectAmountAPI", s)
	controller.Post("/human-suspect-amount:submit", "Submit")
	controller.Post("/human-suspect-amount", "Create")
	controller.Put("/human-suspect-amount", "Update")
	controller.Delete("/human-suspect-amount/{id}", "DeleteById")
	controller.GetData("/human-suspect-amount/{id}", "FindById")
	controller.GetPaging("/human-suspect-amount", "FindPaging")
	return controller
}

func (s *HumanSuspectAmountAPI) Submit(ctx context.Context, cmd *command.HumanSuspectAmountSubmitCommand) (*model.HumanSuspectAmount, error) {
	return &cmd.Data, s.service.Submit(ctx, &cmd.Data)
}

func (s *HumanSuspectAmountAPI) Create(ctx context.Context, cmd *command.HumanSuspectAmountCreateCommand) (*model.HumanSuspectAmount, error) {
	return &cmd.Data, s.service.Create(ctx, &cmd.Data)
}

func (s *HumanSuspectAmountAPI) Update(ctx context.Context, cmd *command.HumanSuspectAmountUpdateCommand) (*model.HumanSuspectAmount, error) {
	return &cmd.Data, s.service.Update(ctx, &cmd.Data)
}

func (s *HumanSuspectAmountAPI) DeleteById(ctx context.Context, qry *query.HumanSuspectAmountFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanSuspectAmountAPI) FindById(ctx context.Context, qry *query.HumanSuspectAmountFindByIdQuery) (*model.HumanSuspectAmount, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *HumanSuspectAmountAPI) FindPaging(ctx context.Context, qry *query.HumanSuspectAmountFindPagingQuery) (store.FindPagingResult[*model.HumanSuspectAmount], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}
