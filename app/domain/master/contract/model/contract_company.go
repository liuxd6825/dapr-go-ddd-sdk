package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// ContractCompany 合同-公司关联
type ContractCompany struct {
	xbase.BaseModel `bson:",inline"`

	ContractId   string `json:"contractId" gorm:"contract_id" bson:"contract_id" index:"" title:"合同ID"`
	RelationType string `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string `json:"name" gorm:"name" bson:"name" title:"公司名称"`
	LegalPerson  string `json:"legalPerson" gorm:"legal_person" bson:"legal_person" title:"法人"`
	Phone        string `json:"phone" gorm:"phone" bson:"phone" title:"预留电话"`
	IdentNum     string `json:"identNum" gorm:"ident_num" bson:"ident_num" title:"纳税人识别号"`
	Addr         string `json:"addr" gorm:"addr" bson:"addr" title:"注册地址"`
}

func NewContractCompany() *ContractCompany {
	return &ContractCompany{}
}