package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// CompanyAccount 公司账号关联
type CompanyAccount struct {
	xbase.BaseModel `bson:",inline"`

	CompanyId    string     `json:"companyId" gorm:"company_id" bson:"company_id" index:"" title:"公司ID"`
	RelationType string     `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	OrgName      string     `json:"orgName" gorm:"org_name" bson:"org_name" title:"机构名称"`
	AccountNum   string     `json:"accountNum" gorm:"account_num" bson:"account_num" title:"账号"`
	StartDate    *time.Time `json:"startDate" gorm:"start_date" bson:"start_date" title:"开始日期"`
	EndDate      *time.Time `json:"endDate" gorm:"end_date" bson:"end_date" title:"结束日期"`
	Used         bool       `json:"used" gorm:"used" bson:"used" title:"是否常用"`
}

func NewCompanyAccount() *CompanyAccount {
	return &CompanyAccount{}
}
