package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/task/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/template/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
)

type Task struct {
	Base      `bson:",inline"`
	DocId     string                 `json:"docId,omitempty" gorm:"doc_id" bson:"doc_id" index:"" desc:"文档id"`
	FileId    string                 `json:"fileId,omitempty" gorm:"file_id" bson:"file_id" index:"" desc:"文件id"`
	FileName  string                 `json:"fileName,omitempty" gorm:"file_name"  bson:"file_name" index:""  desc:"文件名称"`
	SheetName string                 `json:"sheetName,omitempty"  gorm:"sheet_name" bson:"sheet_name"  desc:"Sheet页"`
	StartTime *times.Time            `json:"startTime,omitempty" gorm:"start_time"  bson:"start_time"  desc:"开始时间"`
	EndTime   *times.Time            `json:"endTime,omitempty" gorm:"end_time"  bson:"end_time"  desc:"最后日期"`
	Total     int64                  `json:"total,omitempty" gorm:"total"  bson:"total"  desc:"总记录数"`
	Complete  int64                  `json:"complete,omitempty" gorm:"complete"  bson:"complete" desc:"完成数量"`
	SchemaId  string                 `json:"schemaId,omitempty" gorm:"schema_id"  bson:"schema_id" desc:"主数据类型"`
	MapHeads  []*model.MapHead       `json:"mapHeads" gorm:"map_heads;type:json"  bson:"map_heads"`
	Fields    []*model.TemplateField `json:"fields,omitempty" gorm:"fields;type:json"  bson:"fields" desc:"模板"`
	State     enums.TaskState        `json:"state,omitempty" gorm:"state"  bson:"state" desc:"状态"`
	Message   string                 `json:"message,omitempty" gorm:"message"  bson:"message" desc:"信息"`
	Lock      bool                   `json:"lock" gorm:"lock"  bson:"lock" desc:"锁定"`
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
