package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CreateFile(t *testing.T) {
	xtest.InitEnv_MongoRemoteMaster()
	ctx := xtest.NewContext()
	cmd := &command.ExcelCreateByFileCommand{
		CaseId:   "1001",
		DocId:    "F3YOfd110hzRCKuT7eZPSHL4k",
		FileId:   "FIcuNxpulmMopISDYhPdHyrAK",
		FileName: "赵蕾流水汇总.xlsx",
	}

	service := NewExcelService()
	excelFile, err := service.CreateFile(ctx, cmd)
	assert.NoError(t, err)
	t.Log(excelFile)

	recordService := NewRecordService()
	recordService.Create4Excel(ctx, &command.RecordCreate4ExcelCommand{
		CommandId: cmd.FileId,
		Data: field.RecordCreate4ExcelCommandFields{
			CaseId:    cmd.CaseId,
			DocId:     cmd.DocId,
			FileId:    cmd.FileId,
			FileName:  cmd.FileName,
			SheetName: "1w",
			Template:  NewRecordTemplate10w(),
		},
	}, nil)
}
