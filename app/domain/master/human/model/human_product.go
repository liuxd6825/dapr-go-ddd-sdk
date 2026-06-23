package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// HumanProduct 人员-产品关联
type HumanProduct struct {
	xbase.BaseModel `bson:",inline"`

	HumanId      string `json:"humanId" gorm:"human_id" bson:"human_id" index:"" title:"人员ID"`
	RelationType string `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string `json:"name" gorm:"name" bson:"name" title:"产品名称"`
	Code         string `json:"code" gorm:"code" bson:"code" title:"编码"`
	Nature       string `json:"nature" gorm:"nature" bson:"nature" title:"性质"`
	ProductType  string `json:"productType" gorm:"product_type" bson:"product_type" title:"产品类型"`
	Seller       string `json:"seller" gorm:"seller" bson:"seller" title:"销售方"`
	Remark       string `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewHumanProduct() *HumanProduct {
	return &HumanProduct{}
}