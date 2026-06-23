package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// HumanAccount 人员-账户关联
type HumanAccount struct {
	xbase.BaseModel `bson:",inline"`

	HumanId     string     `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	RelationType string    `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	OrgName     string     `json:"orgName" gorm:"org_name" bson:"org_name" title:"机构名称"`
	AccountNum  string     `json:"accountNum" gorm:"account_num" bson:"account_num" title:"账号"`
	StartDate   *time.Time `json:"startDate" gorm:"start_date" bson:"start_date" title:"起始日期"`
	EndDate     *time.Time `json:"endDate" gorm:"end_date" bson:"end_date" title:"结束日期"`
	Used        bool       `json:"used" gorm:"used" bson:"used" title:"是否使用"`
	Remark      string     `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanAccount() *HumanAccount {
	return &HumanAccount{}
}