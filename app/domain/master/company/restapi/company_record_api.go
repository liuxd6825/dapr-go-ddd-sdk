package restapi

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/company/service"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type CompanyRecordAPI struct {
	service  *service.CompanyRecordService
	rootPath string
}

func NewCompanyRecordAPI(rootPath string) *CompanyRecordAPI {
	return &CompanyRecordAPI{
		rootPath: rootPath,
		service:  service.NewCompanyRecordService(),
	}
}

func (s *CompanyRecordAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.company.CompanyRecordAPI", s)
	controller.GetPaging("/company-record", "FindPaging")
	controller.GetData("/company/{companyId}/company-record", "FindByCompanyId")
	controller.GetOne("/company-record/{id}", "FindById")
	return controller
}

func (s *CompanyRecordAPI) FindById(ctx context.Context, qry *query.CompanyRecordFindByIdQuery) (*model.CompanyRecord, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *CompanyRecordAPI) FindPaging(ctx context.Context, qry *query.CompanyRecordFindPagingQuery) (store.FindPagingResult[*model.CompanyRecord], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *CompanyRecordAPI) FindByCompanyId(ctx context.Context, qry *query.CompanyRecordFindByCompanyIdQuery) (store.FindPagingResult[*model.CompanyRecord], error) {
	res := s.service.FindByCompanyId(ctx, qry, qry.CompanyId)
	return res, res.GetError()
}
