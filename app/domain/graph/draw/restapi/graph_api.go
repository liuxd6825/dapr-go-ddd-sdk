package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type GraphAPI struct {
	env          *env.Env
	graphService *service.GraphService
	rootPath     string
}

func NewGraphAPI(env *env.Env, rootPath string) *GraphAPI {
	graphService := service.NewGraphService()
	return &GraphAPI{
		env:          env,
		graphService: graphService,
		rootPath:     rootPath,
	}
}

func (s *GraphAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/case/{caseId}/draw/{drawId}/graph", "FindByDrawId")
}

func (s *GraphAPI) InitController(app *iris.Application) error {
	controller := restapi.NewController(app, s.rootPath, s)
	controller.GetOne("/case/{caseId}/draw/{drawId}/graph", "FindByDrawId")
	return nil
}

type FindByDrawIdParams struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
	DrawId string `json:"drawId" path:"drawId" required:"true"`
}

func (s *GraphAPI) FindByDrawId(ctx context.Context, params *FindByDrawIdParams) (any, error) {
	graphView := s.graphService.FindById(ctx, params.CaseId, params.DrawId)
	result := response.NewResultList(graphView)
	return result, nil
}
