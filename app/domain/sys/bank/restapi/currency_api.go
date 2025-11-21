package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/bank/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type BankAPI struct {
	env         *env.Env
	BankService *service.BankService
	rootPath    string
}

func NewBankAPI(env *env.Env, rootPath string) *BankAPI {
	return &BankAPI{
		env:         env,
		BankService: service.NewBankService(),
		rootPath:    rootPath,
	}
}

func (s *BankAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath+"/master", "sys.BankApi", s)
	ctl.Post("/bank", "Create")
	ctl.Put("/bank", "Update")
	ctl.Delete("/bank", "Delete", restapi.WithParamsInBody(true))
	ctl.Delete("/bank:deleteBatch", "DeleteBatch", restapi.WithParamsInBody(true))
	ctl.GetOne("/bank/{id}", "FindById")
	ctl.GetPaging("/bank", "FindPaging")
	return ctl
}

func (s *BankAPI) Create(ctx context.Context, cmd *command.BankCreateCommand) error {
	err := tx.StartTx(ctx, []string{s.BankService.GetConfig().DBKey}, func(ctx context.Context, options ...*store.SessionOptions) error {
		return s.BankService.Create(ctx, cmd)
	})
	return err
}

func (s *BankAPI) Update(ctx context.Context, cmd *command.BankUpdateCommand) error {
	return s.BankService.Update(ctx, cmd)
}

func (s *BankAPI) Delete(ctx context.Context, cmd *command.BankDeleteCommand) error {
	return s.BankService.Delete(ctx, cmd)
}

func (s *BankAPI) DeleteBatch(ctx context.Context, cmd *command.BankDeleteBatchCommand) error {
	return s.BankService.DeleteBatch(ctx, cmd)
}

func (s *BankAPI) FindById(ctx context.Context, qry *query.FindByIdQuery) (*model.Bank, error) {
	return s.BankService.FindById(ctx, qry)
}

func (s *BankAPI) FindPaging(ctx context.Context, qry *idao.FindPagingQueryRequest) (idao.FindPagingResult[*model.Bank], error) {
	return s.BankService.FindPaging(ctx, qry)
}
