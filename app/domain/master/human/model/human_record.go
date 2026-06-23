package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// HumanRecord 人员-记录
type HumanRecord struct {
	xbase.BaseModel `bson:",inline"`

	HumanId    string     `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	Title      string     `json:"title" gorm:"title" bson:"title" title:"标题"`
	Content    string     `json:"content" gorm:"content" bson:"content" title:"内容"`
	Status     string     `json:"status" gorm:"status" bson:"status" title:"状态"`
	RecordTime *time.Time `json:"recordTime" gorm:"record_time" bson:"record_time" title:"记录时间"`
	Recorder   string     `json:"recorder" gorm:"recorder" bson:"recorder" title:"记录人"`
	Remark     string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanRecord() *HumanRecord {
	return &HumanRecord{}
}
