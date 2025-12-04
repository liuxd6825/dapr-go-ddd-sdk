package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordAPI struct {
	service  *service.RecordService
	rootPath string
}

func NewRecordAPI(rootPath string) *RecordAPI {
	return &RecordAPI{
		rootPath: rootPath,
		service:  service.NewRecordService(),
	}
}

func (s *RecordAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/master", "master.RecordAPI", s)
	controller.GetOne("/record/{id}", "FindById", restapi.NewAPIOptions())
	controller.GetPaging("/record", "FindPaging")
	controller.GetData("/record:account-name", "AccountFindByName")
	return controller
}

func (s *RecordAPI) FindById(ctx context.Context, qry *query.RecordFindByIdQuery) (any, error) {
	return s.service.FindById(ctx, qry)
}

func (s *RecordAPI) FindPaging(ctx context.Context, qry *query.RecordFindByCaseIdQuery) (store.FindPagingResult[*model.Record], error) {
	res := s.service.FindPagingByCaseId(ctx, qry)
	return res, res.GetError()
}

func (s *RecordAPI) AccountFindByName(ctx context.Context, qry *query.RecordAccountFindByName) ([]*query.RecordAccountFindByNameResult, error) {
	return s.service.AccountFindByName(ctx, qry)
}
