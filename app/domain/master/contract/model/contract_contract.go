package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// ContractContract 合同-合同关联
type ContractContract struct {
	xbase.BaseModel `bson:",inline"`

	ContractId   string  `json:"contractId" gorm:"contract_id" bson:"contract_id" index:"" title:"合同ID"`
	RelationType string  `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string  `json:"name" gorm:"name" bson:"name" title:"合同名称"`
	Code         string  `json:"code" gorm:"code" bson:"code" title:"合同编号"`
	ContractType string  `json:"contractType" gorm:"contract_type" bson:"contract_type" title:"合同类型"`
	FirstParty   string  `json:"firstParty" gorm:"first_party" bson:"first_party" title:"甲方"`
	SecondParty  string  `json:"secondParty" gorm:"second_party" bson:"second_party" title:"乙方"`
	Amount       float64 `json:"amount" gorm:"amount" bson:"amount" title:"合同金额"`
}

func NewContractContract() *ContractContract {
	return &ContractContract{}
}