package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/command"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	task2 "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetImportProgress(t *testing.T) {

}

func TestTaskDomainCmdService_Create(t *testing.T) {
	xtest.InitEnv_MongoRemoteTest()
	ctx := xtest.NewContext()
	// 更新导入进度
	id := idutils.NewId()
	cmd := &command.TaskCreateCommand{}
	cmd.CommandId = id
	cmd.Data = field.TaskCreateFields{
		Id:        id,
		CaseId:    id,
		Name:      id,
		DocId:     id,
		FileId:    id,
		SheetName: id,
	}

	service := NewTaskService()
	err := service.Create(ctx, cmd)
	assert.NoError(t, err)

	t.Logf("Id=%v", id)
	task, err := service.FindById(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, task)

	// 更新导入进度
	now := times.PNow()
	p := &field.TaskUpdateProgressFields{
		Id:           id,
		CompleteRows: 1,
		State:        task2.TaskStateError,
		StartTime:    now,
		EndTime:      now,
		Message:      "message",
	}

	err = service.UpdateProgress(ctx, p)
	assert.NoError(t, err)
}
