package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// HumanCompany 人员-公司关联
type HumanCompany struct {
	xbase.BaseModel `bson:",inline"`

	HumanId      string `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	RelationType string `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string `json:"name" gorm:"name" bson:"name" title:"公司名称"`
	LegalPerson  string `json:"legalPerson" gorm:"legal_person" bson:"legal_person" title:"法定代表人"`
	Phone        string `json:"phone" gorm:"phone" bson:"phone" title:"电话"`
	Addr         string `json:"addr" gorm:"addr" bson:"addr" title:"地址"`
	IdentNum     string `json:"identNum" gorm:"ident_num" bson:"ident_num" title:"证件号"`
	Remark       string `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanCompany() *HumanCompany {
	return &HumanCompany{}
}