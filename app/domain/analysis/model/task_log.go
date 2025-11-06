package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
)

type TaskLog struct {
	xbase.BaseModel
	Time    time.Time `json:"time" gorm:"time" bson:"time" `
	Index   int64     `json:"index" gorm:"index" bson:"index" validate:"required"`
	TaskId  string    `json:"task_id" gorm:"task_id" bson:"task_id" validate:"required"`
	Message string    `json:"message" gorm:"message" bson:"message" validate:"required"`
}

func NewTaskLog(taskId string, index int64, msg string) *TaskLog {
	return &TaskLog{
		BaseModel: xbase.BaseModel{
			Id: idutils.NewId(),
		},
		Time:    time.Now(),
		TaskId:  taskId,
		Index:   index,
		Message: msg,
	}
}
