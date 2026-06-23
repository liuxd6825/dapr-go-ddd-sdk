package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/service"
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
	controller.GetData("/record:distinct-account", "DistinctAccount")
	controller.GetData("/record:distinct-name", "DistinctName")
	controller.GetData("/record:distinct-company", "DistinctCompany")
	controller.GetData("/record:distinct-human", "DistinctHuman")
	return controller
}

func (s *RecordAPI) FindById(ctx context.Context, qry *query.RecordFindByIdQuery) (any, error) {
	return s.service.FindById(ctx, qry)
}

func (s *RecordAPI) FindPaging(ctx context.Context, qry *query.RecordFindByCaseIdQuery) (store.FindPagingResult[*model.Record], error) {
	res := s.service.FindPagingByCaseId(ctx, qry)
	return res, res.GetError()
}

func (s *RecordAPI) DistinctAccount(ctx context.Context, qry *query.DistinctAccountByNameQuery) ([]*query.DistinctAccountResult, error) {
	return s.service.DistinctAccount(ctx, qry)
}

func (s *RecordAPI) DistinctName(ctx context.Context, qry *query.DistinctNameQuery) ([]*query.DistinctNameResult, error) {
	return s.service.DistinctName(ctx, qry)
}

func (s *RecordAPI) DistinctCompany(ctx context.Context, qry *query.DistinctNameQuery) ([]*query.DistinctNameResult, error) {
	return s.service.DistinctCompany(ctx, qry)
}

func (s *RecordAPI) DistinctHuman(ctx context.Context, qry *query.DistinctNameQuery) ([]*query.DistinctNameResult, error) {
	return s.service.DistinctHuman(ctx, qry)
}
