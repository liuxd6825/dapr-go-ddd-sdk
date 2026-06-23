package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// HumanAddress 人员-地址关联
type HumanAddress struct {
	xbase.BaseModel `bson:",inline"`

	HumanId  string     `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	Country  string     `json:"country" gorm:"country" bson:"country" title:"国家"`
	Province string     `json:"province" gorm:"province" bson:"province" title:"省"`
	City     string     `json:"city" gorm:"city" bson:"city" title:"市"`
	Addr     string     `json:"addr" gorm:"addr" bson:"addr" title:"详细地址"`
	Purpose  string     `json:"purpose" gorm:"purpose" bson:"purpose" title:"用途"`
	StartDate *time.Time `json:"startDate" gorm:"start_date" bson:"start_date" title:"起始日期"`
	EndDate   *time.Time `json:"endDate" gorm:"end_date" bson:"end_date" title:"结束日期"`
	Used     bool       `json:"used" gorm:"used" bson:"used" title:"是否使用"`
	Remark   string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanAddress() *HumanAddress {
	return &HumanAddress{}
}