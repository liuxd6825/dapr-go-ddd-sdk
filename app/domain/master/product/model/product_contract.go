package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// ProductContract 产品-合同关联
type ProductContract struct {
	xbase.BaseModel `bson:",inline"`

	ProductId    string  `json:"productId" gorm:"product_id" bson:"product_id" index:"" title:"产品ID"`
	RelationType string  `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string  `json:"name" gorm:"name" bson:"name" title:"合同名称"`
	Code         string  `json:"code" gorm:"code" bson:"code" title:"编码"`
	FirstParty   string  `json:"firstParty" gorm:"first_party" bson:"first_party" title:"甲方"`
	SecondParty  string  `json:"secondParty" gorm:"second_party" bson:"second_party" title:"乙方"`
	Amount       float64 `json:"amount" gorm:"amount" bson:"amount" title:"金额"`
	Remark       string  `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewProductContract() *ProductContract {
	return &ProductContract{}
}