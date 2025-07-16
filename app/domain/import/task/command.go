package task

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

// TaskCreateCommand
// @Description:
type TaskCreateCommand = xbase.Command[TaskCreateFields]
type TaskLockCommand = xbase.Command[TaskLockField]
type TaskDeleteCommand = xbase.DeleteByIdCommand

type TaskRecordCreateCommand = xbase.Command[TaskRecordCreateFields]
type TaskRecordDeleteCommand = xbase.Command[TaskRecordDeleteFields]
type TaskRecordImportCommand = xbase.Command[TaskRecordImportFields]
type TaskRecordNoImportCommand = xbase.Command[TaskRecordNoImportFields]

type TaskStopCommand = xbase.Command[TaskStopFields]
type TaskUpdateStateCommand = xbase.Command[TaskUpdateStateFields]
type TaskUnlockCommand = xbase.Command[TaskLockField]
type TaskUpdateCommand = xbase.Command[TaskUpdateFields]
type TaskUpdateProgressCommand = xbase.Command[TaskUpdateProgressFields]
type TaskUpdateRecordTemplateCommand = xbase.Command[TaskUpdateRecordTemplateFields]
type TaskValidateCommand = xbase.Command[TaskStopFields]
