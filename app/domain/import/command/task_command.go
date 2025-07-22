package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

// TaskCreateCommand
// @Description:
type TaskCreateCommand struct {
	xbase.Command[field.TaskCreateFields]
}
type TaskLockCommand struct {
	xbase.Command[field.TaskLockField]
}
type TaskDeleteCommand struct{ xbase.DeleteByIdCommand }

type TaskRecordCreateCommand struct {
	xbase.Command[field.TaskRecordCreateFields]
}
type TaskRecordDeleteCommand struct {
	xbase.Command[field.TaskRecordDeleteFields]
}
type TaskRecordImportCommand struct {
	xbase.Command[field.TaskRecordImportFields]
}
type TaskRecordNoImportCommand struct {
	xbase.Command[field.TaskRecordNoImportFields]
}

type TaskStopCommand struct {
	xbase.Command[field.TaskStopFields]
}
type TaskUpdateStateCommand struct {
	xbase.Command[field.TaskUpdateStateFields]
}
type TaskUnlockCommand struct {
	xbase.Command[field.TaskLockField]
}
type TaskUpdateCommand struct {
	xbase.Command[field.TaskUpdateFields]
}
type TaskUpdateProgressCommand struct {
	xbase.Command[field.TaskUpdateProgressFields]
}
type TaskUpdateRecordTemplateCommand struct {
	xbase.Command[field.TaskUpdateRecordTemplateFields]
}
type TaskValidateCommand struct {
	xbase.Command[field.TaskStopFields]
}

func NewTaskUpdateProgressCommand(commandId, taskId string) *TaskUpdateProgressCommand {
	cmd := &TaskUpdateProgressCommand{}
	cmd.CommandId = commandId
	cmd.Data = field.TaskUpdateProgressFields{
		Id: taskId,
	}
	return cmd
}
