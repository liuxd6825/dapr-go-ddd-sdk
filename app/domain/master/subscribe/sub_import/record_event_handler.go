package sub_import

import (
	"context"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/factory"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/view"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/tx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/restapi"
)

// RecordEventSubHandler
// @Description: 处理数据导入事件
type RecordEventSubHandler struct {
	rootPath             string
	env                  *env.Env
	factory              *factory.RecordFactory
	recordService        *service.RecordService
	recordDayViewService *service.RecordDayViewService
	tranDetailService    *service.TranDetailService
}

func NewRecordEventHandler(env *env.Env, baseUrl string) *RecordEventSubHandler {
	return &RecordEventSubHandler{
		rootPath:             baseUrl,
		env:                  env,
		factory:              factory.NewRecordFactory(),
		recordService:        service.NewRecordService(),
		recordDayViewService: service.NewRecordDayViewService(),
		tranDetailService:    service.NewTranDetailService(),
	}
}

func (s *RecordEventSubHandler) NewAPIController(app *iris.Application) *restapi.ApiController {
	ctl := restapi.NewController(app, "subscribe/master/event", "RecordEventSubHandler", s)
	ctl.EventHandle("record-import-master-event", "RecordImportMasterEvent")
	ctl.Handle(iris.MethodOptions, "record-import-master-event", "Check")
	return ctl
}

func (s *RecordEventSubHandler) Check(ctx context.Context) error {
	logs.Infofmt(ctx, "graph/master/subscribe/import/record-import-master-event:check")
	return nil
}

func (s *RecordEventSubHandler) RecordImportMasterEvent(ctx context.Context, event *event.RecordImportMasterEvent) error {
	logs.Infofmt(ctx, "record-import-master-event eventId:%s; occurredOn:%s; ", event.Id, event.CreatedTime.Format(time.DateTime))
	records, err := s.factory.NewByRecordImportMasterEvent(ctx, event)
	if err != nil {
		return err
	}
	/*	var details []*model.Tran
		var recordDays []*view.RecordDayView
		for _, record := range records {
			tran := model.NewTranFromRecord(record)
			details = append(details, tran)

			d1, err := view.NewRecordDayView(record)
			if err != nil {
				return err
			}
			recordDays = append(recordDays, d1)
			// recordDays = append(recordDays, d2)
		}
	*/

	return tx.StartTx(ctx, tx.NewTxCfg(config.DBKey), func(ctx context.Context, options ...*store.SessionOptions) error {
		err = s.recordService.CreateMany(ctx, records)
		return err

		/*		err = s.recordDayViewService.IncAmountMany(ctx, recordDays)
				if err != nil {
					return err
				}
				return s.tranDetailService.CreateMany(ctx, details)
		*/
	})

}

func (s *RecordEventSubHandler) newRecordDayView(tran *model.Tran, record *model.Record) *view.RecordDayView {
	v := &view.RecordDayView{
		Id:         tran.Id,
		MasterType: tran.MasterType,
		MasterId:   tran.MasterId,
		TenantId:   tran.TenantId,
		CaseId:     tran.CaseId,
		Name:       tran.Name,
		Acct:       tran.Acct,
		OppName:    tran.OppName,
		OppAcct:    tran.OppAcct,
		Date:       tran.Date,
		Year:       tran.Year,
		Month:      tran.Month,
		Day:        tran.Day,
		Ccy:        tran.Ccy,
		Payout:     record.Payout,
		Income:     record.Income,
		Amount:     tran.Amount,
	}
	return v
}
