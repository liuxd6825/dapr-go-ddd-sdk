package service

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/service/command"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
	"github.com/stretchr/testify/assert"
)

const fileId = "FIcuNxpulmMopISDYhPdHyrAK"
const caseId = "1001"
const docId = "F3YOfd110hzRCKuT7eZPSHL4k"
const taskId = "FIcuNxpulmMopISDYhPdHyrAK"
const fileName = "赵蕾流水汇总.xlsx"

func Test_CreateFile(t *testing.T) {
	xtest2.InitEnv_MongoRemoteMaster()
	ctx := xtest2.NewContext()
	cmd := &command.ExcelFileCreateCommand{}
	cmd.CommandId = "FIcuNxpulmMopISDYhPdHyrAK"
	cmd.Data = command.ExcelCreateByFileCommandData{
		Id:        fileId,
		CaseId:    "1001",
		DocId:     "F3YOfd110hzRCKuT7eZPSHL4k",
		DocFileId: fileId,
		FileName:  "赵蕾流水汇总.xlsx",
	}

	service := NewExcelFileService()
	excelFile, err := service.Create(ctx, cmd)
	assert.NoError(t, err)
	t.Log(excelFile)
}
