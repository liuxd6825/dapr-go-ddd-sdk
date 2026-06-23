package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// CompanyRecord 公司-流水关联（只读视图）
type CompanyRecord struct {
	xbase.BaseModel `bson:",inline"`

	CompanyId    string     `json:"companyId" gorm:"company_id" bson:"company_id" index:"" title:"公司ID"`
	RelationType string     `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Title        string     `json:"title" gorm:"title" bson:"title" title:"标题"`
	Content      string     `json:"content" gorm:"content" bson:"content" title:"内容"`
	Status       string     `json:"status" gorm:"status" bson:"status" title:"分析状态"`
	RecordTime   *time.Time `json:"recordTime" gorm:"record_time" bson:"record_time" title:"记录时间"`
	Recorder     string     `json:"recorder" gorm:"recorder" bson:"recorder" title:"记录人"`
}

func NewCompanyRecord() *CompanyRecord {
	return &CompanyRecord{}
}
