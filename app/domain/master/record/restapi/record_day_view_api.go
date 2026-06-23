package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/record/view"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordDayViewApi struct {
	service  *service.RecordDayViewService
	rootPath string
}

func NewRecordDayViewApi(rootPath string) *RecordDayViewApi {
	return &RecordDayViewApi{
		rootPath: rootPath,
		service:  service.NewRecordDayViewService(),
	}
}

func (s *RecordDayViewApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/master", "master.RecordDayAPI", s)
	controller.GetOne("/record-day/{id}", "FindById")
	controller.GetPaging("/record-day", "FindPaging")
	controller.GetData("/record-day:sum-chart", "FindBySumChart")
	controller.GetData("/record-day:sum-table", "FindBySumTable")
	return controller
}

func (s *RecordDayViewApi) FindPaging(ctx context.Context, qry *query.RecordFindByCaseIdQuery) (store.FindPagingResult[*view.RecordDayView], error) {
	res := s.service.FindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *RecordDayViewApi) FindById(ctx context.Context, qry *query.RecordDayFindByIdQuery) (*view.RecordDayView, error) {
	return s.service.FindById(ctx, qry.Id)
}

func (s *RecordDayViewApi) FindBySumChart(ctx context.Context, qry *query.RecordDayFindBySumChartQuery) ([]*view.RecordSumChartView, error) {
	tenantId := appctx.GetTenantId2(ctx)
	data, _, err := s.service.FindRecordSumChart(ctx, tenantId, qry.CaseId, qry.SummaryType, qry.Filter)
	return data, err
}

func (s *RecordDayViewApi) FindBySumTable(ctx context.Context, qry *query.RecordDayFindBySumTableQuery) (*view.RecordSumTableQueryView, error) {
	tenantId := appctx.GetTenantId2(ctx)
	data, _, err := s.service.FindRecordSumTable(ctx, tenantId, qry.CaseId, qry.Filter, qry.GroupFilter, qry.Sort, qry.PageNum, qry.PageSize)
	return data, err
}
