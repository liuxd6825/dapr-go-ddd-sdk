package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	graph2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type GraphAPI struct {
	env          *env.Env
	rootPath     string
	queryService *graph2.QueryService
}

func NewGraphAPI(env *env.Env, rootPath string) *GraphAPI {
	queryService := graph2.NewQueryService()
	return &GraphAPI{
		env:          env,
		rootPath:     rootPath,
		queryService: queryService,
	}
}

func (s *GraphAPI) InitController(app *iris.Application) error {
	controller := restapi.NewController(app, s.rootPath, s)
	controller.GetOne("/case/{caseId}/master/graph?id={id}", "FindById")
	controller.GetList("/case/{caseId}/master/graph", "FindByCaseId")
	return nil
}

type FindByCaseIdParams struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
}

func (s *GraphAPI) FindByCaseId(ctx context.Context, params *FindByCaseIdParams) (any, error) {
	graphView := s.queryService.FindByCaseId(ctx, params.CaseId)
	resData := response.NewResultList(graphView)
	return resData, nil
}

type FindByIdParams struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
	Id     string `json:"id" query:"id" required:"true"`
}

func (s *GraphAPI) FindById(ctx context.Context, params *FindByIdParams) (any, error) {
	graphView := s.queryService.FindById(ctx, params.CaseId, params.Id)
	data := response.NewResultList(graphView)
	return data, nil
}

type FindByNameParams struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
	Name   string `json:"name" path:"name" required:"true"`
}

func (s *GraphAPI) FindByName(ctx context.Context, params *FindByNameParams) (any, error) {
	graphView := s.queryService.FindById(ctx, params.CaseId, params.Name)
	resData := response.NewResultList(graphView)
	return resData, nil
}
