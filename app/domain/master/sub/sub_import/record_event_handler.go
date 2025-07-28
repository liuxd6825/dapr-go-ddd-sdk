package sub_import

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordEventSubHandler struct {
	rootPath string
	env      *env.Env
}

func NewRecordEventHandler(env *env.Env, baseUrl string) *RecordEventSubHandler {
	return &RecordEventSubHandler{
		rootPath: baseUrl,
		env:      env,
	}
}

func (s *RecordEventSubHandler) InitController(app *iris.Application) error {
	ctl := restapi.NewController(app, "/subscribe/import", s)
	ctl.EventHandle("record-import-master-event", "RecordImportMasterEvent")
	ctl.Handle(iris.MethodOptions, "record-import-master-event", "Check")
	return nil
}

func (s *RecordEventSubHandler) Check(ctx context.Context) error {
	println("subscribe/import/record-import-master-event:check")
	return nil
}

func (s *RecordEventSubHandler) RecordImportMasterEvent(ctx context.Context, event *event.RecordImportMasterEvent) error {
	println("subscribe/import/record-import-master-event", event)
	return nil
}
