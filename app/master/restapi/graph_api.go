package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	graph2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/master/service/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/master/service/graph/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/irisutils"
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
	b.Handle("POST", "/cdc-mysql", "DataChange")
	b.Handle("POST", "/graph", "Query")
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
		irisutils.SetError(ctx, e)
	})

}

func (s *GraphAPI) Query(ctx iris.Context) {
	gp.Try(func() error {
		data, err := s.queryService.FindById(ctx, "1001", "", "")
		if err != nil {
			return err
		}
		return irisutils.SetData(ctx, data)
	}).Catch(func(e error) {
		irisutils.SetError(ctx, e)
	})
}
