package factory

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/mapperutils"
)

type RecordFactory struct {
}

func NewRecordFactory() *RecordFactory {
	return &RecordFactory{}
}

func (f *RecordFactory) NewByRecordImportMasterEvent(ctx context.Context, e *event.RecordImportMasterEvent) ([]*model.Record, error) {
	list := make([]*model.Record, len(e.Data.Items))
	for i, fields := range e.Data.Items {
		v := model.NewRecord()
		if err := mapperutils.Mapper(fields, v); err != nil {
			return nil, err
		}

		v.TenantId = appctx.GetTenantId2(ctx)
		v.CaseId = e.Data.CaseId
		v.DocId = e.Data.DocId
		v.TaskId = e.Data.TaskId
		v.FileId = e.Data.FileId
		v.SheetId = e.Data.SheetId
		v.RowNum = fields.RowNum
		if fields.Date != nil {
			v.Year = fields.Date.Year()
			v.Month = int(fields.Date.Month())
			v.Day = fields.Date.Day()
		}
		v.TranId = model.NewTranId(v)
		list[i] = v
	}
	return list, nil
}
