package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	model2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

type SuTaskRecordApi struct {
	accountService *service.SuRecordService
	env            *env.Env
	rootPath       string
}

func NewSuTaskRecordApi(env *env.Env, rootPath string) *SuTaskRecordApi {
	accountService := service.NewSuRecordService()
	return &SuTaskRecordApi{
		env:            env,
		rootPath:       rootPath,
		accountService: accountService,
	}
}

func (s *SuTaskRecordApi) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath+"/analysis/", "analysis.SuRecordApi", s)
	controller.GetOne("su-record/{id}", "FindById")
	controller.GetPaging("su-record", "SuRecordFindPaging")
	return controller
}

func (s *SuTaskRecordApi) SuRecordFindPaging(ctx context.Context, qry *query.SuRecordFindByCaseIdQuery) (store.FindPagingResult[*model2.SuRecord], error) {
	res := s.accountService.SuRecordFindPaging(ctx, qry)
	return res, res.GetError()
}

func (s *SuTaskRecordApi) FindById(ctx context.Context, qry *query.SuRecordFindByIdQuery) (*model2.SuRecord, error) {
	task, err := s.accountService.FindById(ctx, qry)
	return task, err
}
