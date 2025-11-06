package service

import (
	"context"
	"testing"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/event"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/events"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

var recordService *RecordService
var ctx context.Context

func init() {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx = xtest2.NewContext()
	recordService = NewRecordService()
}

func TestRecordService_Create(t *testing.T) {
	cmd := &command.RecordCreate4ExcelCommand{}
	cmd.CommandId = "0001"
	cmd.Data = field.RecordCreate4ExcelCommandFields{
		BatchSize: 100,
		CaseId:    "1001",
		DocId:     "F3YOfd110hzRCKuT7eZPSHL4k",
		FileId:    "FIcuNxpulmMopISDYhPdHyrAK",
		FileName:  "10w.xlsx",
		IsView:    false,
		SheetName: "1w",
		TaskId:    taskId,
		Template:  NewRecordTemplate10w(),
	}
	res, err := recordService.Create4Excel(ctx, cmd, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)

}

func TestRecordService_FindPagingByTaskId(t *testing.T) {
	qry := query.NewRecordIeFindPagingByTaskIdQuery("5ZCLHw2Uw1ynsmXYjvAYlsY2", true)
	findRes := recordService.FindPagingByTaskId(ctx, qry)
	if findRes.GetError() != nil {
		t.Error(findRes.GetError())
		return
	}
	t.Log(findRes.GetDataLength())
	t.Log(findRes.GetData())
}

func TestRecordService_Import2Master(t *testing.T) {
	cmd := &command.RecordImport2MasterCommand{}
	cmd.CommandId = idutils.NewUlid2()
	cmd.Data = field.RecordImport2MasterFields{
		CaseId:   "1001",
		DocId:    "F3YOfd110hzRCKuT7eZPSHL4k",
		FileId:   "FIcuNxpulmMopISDYhPdHyrAK",
		FileName: "10w.xlsx",
		TaskId:   "taskId",
		PageSize: 1000,
	}

	err := recordService.Import2Master(ctx, cmd)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRecordService_UpdateFilterCommand(t *testing.T) {
	cmd := &command.RecordUpdateFilterCommand{}
	cmd.CommandId = idutils.NewUlid2()
	cmd.Data = field.RecordIeUpdateFilterFields{
		TaskId: "5ZCLHw2Uw1ynsmXYjvAYlsY2",
		Filter: "(acct=contains='YYYYYYY')",
		Values: map[string]any{"Acct": "11111", "BankName": "222222"},
	}

	err := recordService.UpdateByFilter(ctx, cmd)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestRecordService_UpdateByQueryCommand(t *testing.T) {
	e := &event.RecordImportMasterEvent{}
	e.EventId = idutils.NewUlid2()
	e.OccurredOn = time.Now()
	e.Data = event.RecordImportMasterEventData{
		CaseId:   "1001",
		DocId:    "F3YOfd110hzRCKuT7eZPSHL4k",
		FileName: "10w.xlsx",
		Items: []*field.RecordFields{
			{
				Id:   "1",
				Acct: "111111",
			},
		},
	}

	mapData, err := events.StructToMap(e)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(mapData)
}
