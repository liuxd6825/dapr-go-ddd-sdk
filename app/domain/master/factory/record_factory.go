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
	var list []*model.Record
	for _, fields := range e.Data.Items {
		v := model.NewRecord()
		if err := mapperutils.Mapper(fields, v); err != nil {
			return nil, err
		}
		v.TenantId = appctx.GetTenantId2(ctx)
		v.CaseId = e.Data.CaseId
		v.DocId = e.Data.DocId
		v.TaskId = e.Data.TaskId
		v.FileId = e.Data.FileId
		v.FileName = e.Data.FileName
		list = append(list, v)
	}
	return list, nil
}
