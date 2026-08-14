package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type MysqlQueryAPI struct {
	service  *service.MysqlService
	env      *env.Env
	rootPath string
}

func NewMysqlQueryAPI(env *env.Env, rootPath string) *MysqlQueryAPI {
	return &MysqlQueryAPI{
		service:  service.NewMysqlService(env, rootPath),
		env:      env,
		rootPath: rootPath,
	}
}

func (s *MysqlQueryAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.MysqlQueryAPI", s)
	controller.Post("aggregate", "Aggregate")
	controller.Post("summary", "Summary")
	return controller
}

func (s *MysqlQueryAPI) Aggregate(ctx context.Context, qry *query.MysqlQuery) (any, error) {
	return s.service.Aggregate(ctx, qry)
}

func (s *MysqlQueryAPI) Summary(ctx context.Context, qry *query.SummaryQuery) (any, error) {
	return s.service.Summary(ctx, qry)
}
