package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordChartQueryAPI struct {
	service    *service.RecordChartService
	env        *envs.Env
	rootPath   string
	dorisDbKey string
}

func NewRecordChartQueryAPI(env *envs.Env, rootPath string, dorisDbKey string) *RecordChartQueryAPI {
	return &RecordChartQueryAPI{
		env:        env,
		rootPath:   rootPath,
		dorisDbKey: dorisDbKey,
		service:    service.NewRecordChartService(env, rootPath, dorisDbKey),
	}
}

func (s *RecordChartQueryAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.MysqlQueryAPI", s)
	controller.Post("aggregate", "Aggregate")
	controller.Post("summary", "Summary")
	controller.GetData("autocomplete", "AutoComplete")
	return controller
}

func (s *RecordChartQueryAPI) Aggregate(ctx context.Context, qry *query.MysqlQuery) (any, error) {
	return s.service.Aggregate(ctx, qry)
}

func (s *RecordChartQueryAPI) Summary(ctx context.Context, qry *query.SummaryQuery) (any, error) {
	return s.service.Summary(ctx, qry)
}

func (s *RecordChartQueryAPI) AutoComplete(ctx context.Context, qry *query.AutoCompleteQuery) ([]string, error) {
	return s.service.AutoComplete(ctx, qry)
}
