package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// CompanyHuman 公司-人员关联
type CompanyHuman struct {
	xbase.BaseModel `bson:",inline"`

	CompanyId    string `json:"companyId" gorm:"company_id" bson:"company_id" index:"" title:"公司ID"`
	RelationType string `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string `json:"name" gorm:"name" bson:"name" title:"姓名"`
	Gender       string `json:"gender" gorm:"gender" bson:"gender" title:"性别"`
	Age          int    `json:"age" gorm:"age" bson:"age" title:"年龄"`
	IdentNum     string `json:"identNum" gorm:"ident_num" bson:"ident_num" title:"证件号"`
	Link         string `json:"link" gorm:"link" bson:"link" title:"联系方式"`
	Addr         string `json:"addr" gorm:"addr" bson:"addr" title:"地址"`
}

func NewCompanyHuman() *CompanyHuman {
	return &CompanyHuman{}
}
