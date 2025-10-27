package master

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/cdc"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordHandler struct {
	env           *env.Env
	rootPath      string
	recordService *RecordService
}

func NewRecordHandler(env *env.Env, rootPath string) *RecordHandler {
	return &RecordHandler{
		env:           env,
		rootPath:      rootPath,
		recordService: NewRecordService(),
	}
}

func (s *RecordHandler) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, s.rootPath, "CdcAPI", s)
	ctl.Handle(iris.MethodPost, "/master-cdc-analysis", "DataChange")
	ctl.Handle(iris.MethodOptions, "/master-cdc-analysis", "DaprOptions")
	return ctl
}

func (s *RecordHandler) DaprOptions(ctx context.Context) error {
	return nil
}

func (s *RecordHandler) DataChange(ctx context.Context, cdc *cdc.CDCRecord) error {
	if cdc == nil {
		return nil
	}
	if cdc.Table == "master_record" {

	}
	return nil
}
