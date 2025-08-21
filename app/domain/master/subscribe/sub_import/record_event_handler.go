package sub_import

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/factory"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

type RecordEventSubHandler struct {
	rootPath          string
	env               *env.Env
	factory           *factory.RecordFactory
	recordService     *service.RecordService
	tranDetailService *service.TranDetailService
}

func NewRecordEventHandler(env *env.Env, baseUrl string) *RecordEventSubHandler {
	return &RecordEventSubHandler{
		rootPath:          baseUrl,
		env:               env,
		factory:           factory.NewRecordFactory(),
		recordService:     service.NewRecordService(),
		tranDetailService: service.NewTranDetailService(),
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
	records, err := s.factory.NewByRecordImportMasterEvent(ctx, event)
	if err != nil {
		return err
	}
	var details []*model.Tran
	for _, record := range records {
		details = append(details, model.NewTranFromRecord(record))
	}

	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		err = s.recordService.CreateMany(ctx, records)
		if err != nil {
			return err
		}
		return s.tranDetailService.CreateMany(ctx, details)
	})

}
