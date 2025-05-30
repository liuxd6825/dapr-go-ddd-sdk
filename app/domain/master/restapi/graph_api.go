package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	graph2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/response"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

type GraphAPI struct {
	env          *env.Env
	rootPath     string
	cdcService   *graph2.CdcService
	queryService *graph2.QueryService
}

func NewGraphAPI(env *env.Env, rootPath string) *GraphAPI {
	cdcService := graph2.NewCdcService().Init()
	queryService := graph2.NewQueryService()
	return &GraphAPI{
		env:          env,
		rootPath:     rootPath,
		cdcService:   cdcService,
		queryService: queryService,
	}
}

func (s *GraphAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/cdc-mysql", "DataChange")
	b.Handle(iris.MethodGet, "/case/{caseId}/master/graph?id={id}", "FindById")
	b.Handle(iris.MethodGet, "/case/{caseId}/master/graph?id={id}", "FindById")
}

func (s *GraphAPI) DataChange(ctx iris.Context) {
	gp.Try(func() error {
		var record model.Record

		// 反序列化请求体到结构体
		if err := ctx.ReadJSON(&record); err != nil {
			return err
		}
		logs.InfoMsg(ctx, "record ", "opType=", record.OpType, "table=", record.Table)

		// 根据操作类型处理数据
		switch record.OpType {
		case "c":
			s.cdcService.Create(&record)
		case "u":
			s.cdcService.Update(&record)
		case "d":
			s.cdcService.Delete(&record)
		}
		ctx.StatusCode(iris.StatusOK)
		return nil
	}).Catch(func(e error) {
		web.SetError(ctx, e)
	})

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
