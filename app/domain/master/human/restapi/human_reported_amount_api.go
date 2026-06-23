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

// HumanReportedAmountAPI 报案人金额REST API 控制器
type HumanReportedAmountAPI struct {
	service  *service.HumanReportedAmountService
	rootPath string
}

func NewHumanReportedAmountAPI(rootPath string) *HumanReportedAmountAPI {
	return &HumanReportedAmountAPI{
		rootPath: rootPath,
		service:  service.NewHumanReportedAmountService(),
	}
}

func (s *HumanReportedAmountAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.human.HumanReportedAmountAPI", s)
	controller.Post("/human-reported-amount:submit", "Submit")
	controller.Post("/human-reported-amount", "Create")
	controller.Put("/human-reported-amount", "Update")
	controller.Delete("/human-reported-amount/{id}", "DeleteById")
	controller.GetOne("/human-reported-amount/{id}", "FindById")
	controller.GetPaging("/human-reported-amount", "FindPaging")
	return controller
}

func (s *HumanReportedAmountAPI) Submit(ctx context.Context, cmd *command.HumanReportedAmountSubmitCommand) (*model.HumanReportedAmount, error) {
	return &cmd.Data, s.service.Submit(ctx, &cmd.Data)
}

func (s *HumanReportedAmountAPI) Create(ctx context.Context, cmd *command.HumanReportedAmountCreateCommand) (*model.HumanReportedAmount, error) {
	return &cmd.Data, s.service.Create(ctx, &cmd.Data)
}

func (s *HumanReportedAmountAPI) Update(ctx context.Context, cmd *command.HumanReportedAmountUpdateCommand) (*model.HumanReportedAmount, error) {
	return &cmd.Data, s.service.Update(ctx, &cmd.Data)
}

func (s *HumanReportedAmountAPI) DeleteById(ctx context.Context, qry *query.HumanReportedAmountFindByIdQuery) (any, error) {
	return qry.Id, s.service.DeleteById(ctx, qry.Id)
}

func (s *HumanReportedAmountAPI) FindById(ctx context.Context, qry *query.HumanReportedAmountFindByIdQuery) (*model.HumanReportedAmount, error) {
	amount, err := s.service.FindById(ctx, qry.Id)
	return amount, err
}

func (s *HumanReportedAmountAPI) FindPaging(ctx context.Context, qry *query.HumanReportedAmountFindPagingQuery) (store.FindPagingResult[*model.HumanReportedAmount], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}
