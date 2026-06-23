package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// HumanLink 人员-联系方式
type HumanLink struct {
	xbase.BaseModel `bson:",inline"`

	HumanId     string     `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	LinkType    string     `json:"linkType" gorm:"link_type" bson:"link_type" title:"联系方式类型"`
	LinkContent string     `json:"linkContent" gorm:"link_content" bson:"link_content" title:"联系方式内容"`
	StartDate   *time.Time `json:"startDate" gorm:"start_date" bson:"start_date" title:"起始日期"`
	EndDate     *time.Time `json:"endDate" gorm:"end_date" bson:"end_date" title:"结束日期"`
	Used        bool       `json:"used" gorm:"used" bson:"used" title:"是否使用"`
	Remark      string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanLink() *HumanLink {
	return &HumanLink{}
}