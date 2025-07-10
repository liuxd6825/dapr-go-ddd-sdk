package restapi

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/draw/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
	"strings"
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

func (s *GraphAPI) InitController(app *iris.Application) error {
	h := restapi.NewController(app, s.rootPath, s)
	h.GetData("/case/{caseId}/graph", "FindByName")
	return nil
}

type FindByNameParams struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
	Name   string `json:"name" query:"name" required:"true"`
}

func (s *GraphAPI) FindByName(ctx context.Context, params *FindByNameParams) (any, error) {
	names := strings.Split(params.Name, ",")
	graphView := s.graphService.FindInCaseByNames(ctx, params.CaseId, names)
	result := response.NewResultList(graphView)
	return result, nil
}
