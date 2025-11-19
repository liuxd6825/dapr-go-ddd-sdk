package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/currency/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CurrencyAPI struct {
	env             *env.Env
	currencyService *service.CurrencyService
	rootPath        string
}

func NewCurrencyAPI(env *env.Env, rootPath string) *CurrencyAPI {
	return &CurrencyAPI{
		env:             env,
		currencyService: service.NewCurrencyService(),
		rootPath:        rootPath,
	}
}

func (s *CurrencyAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/master", "sys.CurrencyApi", s)
	ctl.Post("/currency", "Create")
	ctl.Put("/currency", "Update")
	ctl.Delete("/currency", "Delete", restapi.WithParamsInBody(true))
	ctl.Delete("/currency:deleteBatch", "DeleteBatch", restapi.WithParamsInBody(true))
	ctl.GetOne("/currency/{id}", "FindById")
	ctl.GetPaging("/currency", "FindPaging")
	return ctl
}

func (s *CurrencyAPI) Create(ctx context.Context, cmd *command.CurrencyCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.currencyService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		return s.currencyService.Create(ctx, cmd)
	})
	return err
}

func (s *CurrencyAPI) Update(ctx context.Context, cmd *command.CurrencyUpdateCommand) error {
	return s.currencyService.Update(ctx, cmd)
}

func (s *CurrencyAPI) Delete(ctx context.Context, cmd *command.CurrencyDeleteCommand) error {
	return s.currencyService.Delete(ctx, cmd)
}

func (s *CurrencyAPI) DeleteBatch(ctx context.Context, cmd *command.CurrencyDeleteBatchCommand) error {
	return s.currencyService.DeleteBatch(ctx, cmd)
}

func (s *CurrencyAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Currency, error) {
	return s.currencyService.FindById(ctx, qry)
}

func (s *CurrencyAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.Currency], error) {
	return s.currencyService.FindPaging(ctx, qry)
}
