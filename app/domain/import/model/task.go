package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type Base = xbase.BaseModel

type Progress struct {
	Base    `bson:",inline"`
	TaskId  string     `json:"taskId" bson:"task_id" index:""`
	Type    string     `json:"type" bson:"type" index:""`
	Time    times.Time `json:"time" bson:"time" index:""`
	Message string     `json:"message" bson:"message"`
}

type Task struct {
	Base       `bson:",inline"`
	DocId      string           `json:"docId,omitempty" gorm:"doc_id" bson:"doc_id" index:"" title:"文档id"`
	FileId     string           `json:"fileId,omitempty" gorm:"file_id" bson:"file_id" index:"" title:"文件id"`
	FileName   string           `json:"fileName,omitempty" gorm:"file_name"  bson:"file_name" index:""  title:"文件名称"`
	SheetName  string           `json:"sheetName,omitempty"  gorm:"sheet_name" bson:"sheet_name"  title:"Sheet页"`
	StartTime  *times.Time      `json:"startTime,omitempty" gorm:"start_time"  bson:"start_time"  title:"开始时间"`
	EndTime    *times.Time      `json:"endTime,omitempty" gorm:"end_time"  bson:"end_time"  title:"最后日期"`
	Total      int64            `json:"total,omitempty" gorm:"total"  bson:"total"  title:"总记录数"`
	Complete   int64            `json:"complete,omitempty" gorm:"complete"  bson:"complete" title:"完成数量"`
	SchemaId   string           `json:"schemaId,omitempty" gorm:"schema_id"  bson:"schema_id" title:"主数据类型ID"`
	SchemaName string           `json:"schemaName,omitempty" gorm:"schema_nam"  bson:"schema_name" title:"主数据类型"`
	MapHeads   []*MapHead       `json:"mapHeads" gorm:"map_heads;type:json"  bson:"map_heads"`
	Fields     []*TemplateField `json:"fields,omitempty" gorm:"fields;type:json"  bson:"fields" title:"模板"`
	State      TaskState        `json:"state,omitempty" gorm:"state"  bson:"state" title:"状态"`
	Message    string           `json:"message,omitempty" gorm:"message"  bson:"message" title:"信息"`
	Lock       bool             `json:"lock" gorm:"lock"  bson:"lock" title:"锁定"`
}

type TaskState string

const (
	TaskStateNone      TaskState = ""
	TaskStateEditing   TaskState = "编辑中"
	TaskStateGenerated TaskState = "已生成"
	TaskStateImported  TaskState = "已导入"
	TaskStateError     TaskState = "导入错误"
)

func (t TaskState) Name() string {
	return string(t)
}

const (
	TaskAggregateType = "master_import_service.TaskAggregate"
)

func (r *Task) GetAggregateId() string {
	return r.Id
}

func (r *Task) GetAggregateType() string {
	return TaskAggregateType
}

func (r *Task) GetAggregateVersion() string {
	return "v1.0"
}

func (r *Task) GetId() string {
	return r.Id
}

func (r *Task) SetId(v string) {
	r.Id = v
}

func (r *Task) GetTenantId() string {
	return r.TenantId
}

func (r *Task) SetTenantId(v string) {
	r.TenantId = v
}
