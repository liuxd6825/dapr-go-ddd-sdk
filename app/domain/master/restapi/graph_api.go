package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	graph2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
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

func (s *GraphAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodGet, "/case/{caseId}/master/graph?id={id}", "FindById")
	b.Handle(iris.MethodGet, "/case/{caseId}/master/graph", "FindByCaseId")
}

func (s *GraphAPI) FindByCaseId(ctx iris.Context) {
	gp.Try(func() error {
		caseId := ctx.URLParam("caseId")
		graphView := s.queryService.FindByCaseId(ctx, caseId)
		resData := response.NewResultList(graphView)
		return web.SetData(ctx, resData)
	}).Catch(func(e error) {
		web.SetError(ctx, e)
	})
}

func (s *GraphAPI) FindById(ctx iris.Context) {
	gp.Try(func() error {
		caseId := ctx.URLParam("caseId")
		id := ctx.URLParam("id")
		graphView := s.queryService.FindById(ctx, caseId, id)
		resData := response.NewResultList(graphView)
		return web.SetData(ctx, resData)
	}).Catch(func(e error) {
		web.SetError(ctx, e)
	})
}

func (s *GraphAPI) FindByName(ctx iris.Context) {
	gp.Try(func() error {
		caseId := ctx.URLParam("caseId")
		name := ctx.URLParam("name")
		graphView := s.queryService.FindById(ctx, caseId, name)
		resData := response.NewResultList(graphView)
		return web.SetData(ctx, resData)
	}).Catch(func(e error) {
		web.SetError(ctx, e)
	})
}
