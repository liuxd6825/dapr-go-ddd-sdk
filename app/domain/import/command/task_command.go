package command

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/field"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

// TaskCreateCommand
// @Description:
type TaskCreateCommand = xbase.Command[field.TaskCreateFields]
type TaskLockCommand = xbase.Command[field.TaskLockField]
type TaskDeleteCommand = xbase.DeleteByIdCommand

type TaskRecordCreateCommand = xbase.Command[field.TaskRecordCreateFields]
type TaskRecordDeleteCommand = xbase.Command[field.TaskRecordDeleteFields]
type TaskRecordImportCommand = xbase.Command[field.TaskRecordImportFields]
type TaskRecordNoImportCommand = xbase.Command[field.TaskRecordNoImportFields]

type TaskStopCommand = xbase.Command[field.TaskStopFields]
type TaskUpdateStateCommand = xbase.Command[field.TaskUpdateStateFields]
type TaskUnlockCommand = xbase.Command[field.TaskLockField]
type TaskUpdateCommand = xbase.Command[field.TaskUpdateFields]
type TaskUpdateProgressCommand = xbase.Command[field.TaskUpdateProgressFields]
type TaskUpdateRecordTemplateCommand = xbase.Command[field.TaskUpdateRecordTemplateFields]
type TaskValidateCommand = xbase.Command[field.TaskStopFields]

func NewTaskUpdateProgressCommand(commandId, taskId string) *TaskUpdateProgressCommand {
	cmd := &TaskUpdateProgressCommand{}
	cmd.CommandId = commandId
	cmd.Data = field.TaskUpdateProgressFields{
		Id: taskId,
	}
	return cmd
}
