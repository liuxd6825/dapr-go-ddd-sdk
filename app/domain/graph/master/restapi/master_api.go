package restapi

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/graph/vis_network"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type MasterAPI struct {
	env          *env.Env
	rootPath     string
	queryService *service.MasterQueryService
}

func NewMasterAPI(env *env.Env, rootPath string) *MasterAPI {
	queryService := service.NewMasterQueryService()
	return &MasterAPI{
		env:          env,
		queryService: queryService,
		rootPath:     rootPath,
	}
}

type FindByCaseIdRequest struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
}

type FindByIdParams struct {
	CaseId string `json:"caseId" path:"caseId" required:"true"`
	Id     string `json:"id" query:"id" required:"true"`
}

func (s *MasterAPI) NewAPIController(app *iris.Application) *restapi.ApiController {
	controller := restapi.NewController(app, s.rootPath, "master.GraphAPI", s)
	controller.GetOne("/case/{caseId}/master/graph?id={id}", "FindById")
	controller.GetData("/case/{caseId}/master/graph", "FindByCaseId")
	return controller
}

func (s *MasterAPI) FindByCaseId(ctx context.Context, request *FindByCaseIdRequest) *vis_network.ResultList {
	graphView := s.queryService.FindByCaseId(ctx, request.CaseId)
	resData := response.NewResultList(graphView)
	return resData
}

func (s *MasterAPI) FindById(ctx context.Context, params *FindByIdParams) (any, error) {
	graphView := s.queryService.FindById(ctx, params.CaseId, params.Id)
	data := response.NewResultList(graphView)
	return data, nil
}
