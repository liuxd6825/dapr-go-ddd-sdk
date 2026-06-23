package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// CompanyCompany 公司-公司关系
type CompanyCompany struct {
	xbase.BaseModel `bson:",inline"`

	CompanyId    string `json:"companyId" gorm:"company_id" bson:"company_id" index:"" title:"本公司ID"`
	RelationType string `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string `json:"name" gorm:"name" bson:"name" title:"公司名称"`
	LegalPerson  string `json:"legalPerson" gorm:"legal_person" bson:"legal_person" title:"法人"`
	Phone        string `json:"phone" gorm:"phone" bson:"phone" title:"预留电话"`
	IdentNum     string `json:"identNum" gorm:"ident_num" bson:"ident_num" title:"纳税人识别号"`
	Addr         string `json:"addr" gorm:"addr" bson:"addr" title:"注册地址"`
}

func NewCompanyCompany() *CompanyCompany {
	return &CompanyCompany{}
}
