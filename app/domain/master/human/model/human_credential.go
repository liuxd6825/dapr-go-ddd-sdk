package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// HumanCredential 人员-证件
type HumanCredential struct {
	xbase.BaseModel `bson:",inline"`

	HumanId  string     `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	IdentType string    `json:"identType" gorm:"ident_type" bson:"ident_type" title:"证件类型"`
	IdentNum string     `json:"identNum" gorm:"ident_num" bson:"ident_num" title:"证件号"`
	StartDate *time.Time `json:"startDate" gorm:"start_date" bson:"start_date" title:"起始日期"`
	EndDate   *time.Time `json:"endDate" gorm:"end_date" bson:"end_date" title:"结束日期"`
	Used     bool       `json:"used" gorm:"used" bson:"used" title:"是否使用"`
	Remark   string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanCredential() *HumanCredential {
	return &HumanCredential{}
}