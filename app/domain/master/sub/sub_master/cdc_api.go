package sub_master

import (
	"context"
	"github.com/kataras/iris/v12"
	graph2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service/graph/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
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

func (s *CdcAPI) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, s.rootPath, s)
	ctl.Handle(iris.MethodPost, "/master-cdc-graph", "DataChange")
	ctl.Handle(iris.MethodOptions, "/master-cdc-graph", "DaprOptions")
	return nil
}

func (s *CdcAPI) DaprOptions(ctx context.Context) error {
	return nil
}

func (s *CdcAPI) DataChange(ctx context.Context, record *model.Record) error {
	logs.InfoMsg(ctx, "record", " opType=", record.OpType, " table=", record.Table)

	dbSch := s.cdcService.GetDBSchema(record.Table)
	if dbSch == nil {
		return nil
	}
	record.DBSchema = dbSch

	// 根据操作类型处理数据
	switch record.OpType {
	case "c":
		s.cdcService.Create(record)
	case "u":
		s.cdcService.Update(record)
	case "d":
		s.cdcService.Delete(record)
	}
	return nil
}
