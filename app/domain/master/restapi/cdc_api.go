package restapi

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	graph2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/web"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
)

type CdcAPI struct {
	env          *env.Env
	rootPath     string
	cdcService   *graph2.CdcService
	queryService *graph2.QueryService
}

func NewCdcAPI(env *env.Env, rootPath string) *CdcAPI {
	cdcService := graph2.NewCdcService().Init()
	return &CdcAPI{
		env:        env,
		rootPath:   rootPath,
		cdcService: cdcService,
	}
}

func (s *CdcAPI) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(iris.MethodPost, "/master-cdc-graph", "DataChange")
}

func (s *CdcAPI) DataChange(ctx iris.Context) {
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
