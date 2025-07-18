package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

var recordService *RecordService
var ctx context.Context

const taskId = "taskId"

func init() {
	xtest.InitEnv_MongoRemoteMaster()
	ctx = xtest.NewContext()
	recordService = NewRecordService()
}

func TestRecordService_Create(t *testing.T) {
	res, err := recordService.Create4Excel(ctx, &command.RecordCreate4ExcelCommand{
		CommandId: "0001",
		Data: field.RecordCreate4ExcelCommandFields{
			BatchSize: 100,
			CaseId:    "1001",
			DocId:     "F3YOfd110hzRCKuT7eZPSHL4k",
			FileId:    "FIcuNxpulmMopISDYhPdHyrAK",
			FileName:  "10w.xlsx",
			IsView:    false,
			SheetName: "1w",
			TaskId:    taskId,
			Template:  NewRecordTemplate10w(),
		},
	}, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)

}

func TestRecordService_Import2Master(t *testing.T) {
	qry := query.NewRecordIeFindPagingByTaskIdQuery(taskId, true)
	findRes := recordService.FindPagingByTaskId(ctx, qry)
	if findRes.GetError() != nil {
		t.Error(findRes.GetError())
		return
	}
	if findRes.GetDataLength() > 0 {
		t.Error("error record count>0")
		return
	}

	err := recordService.Import2Master(ctx, &command.RecordImport2MasterAppCmd{
		CommandId: "0002",
		Data: field.RecordImport2MasterFields{
			CaseId:   "1001",
			DocId:    "F3YOfd110hzRCKuT7eZPSHL4k",
			FileId:   "FIcuNxpulmMopISDYhPdHyrAK",
			FileName: "10w.xlsx",
			TaskId:   taskId,
			PageSize: 1000,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}
