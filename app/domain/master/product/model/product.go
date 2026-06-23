package model

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

// Product
// @Description: 产品基本信息
type Product struct {
	xbase.BaseModel `bson:",inline"`

	Code               string   `json:"code" gorm:"code" bson:"code" title:"编码"`
	Name               string   `json:"name" gorm:"name" bson:"name" index:"" title:"产品名称"`
	ProductType        string   `json:"productType" gorm:"product_type" bson:"product_type" title:"产品类型"`
	ProductNature      string   `json:"productNature" gorm:"product_nature" bson:"product_nature" title:"产品性质"`
	FundraisingMethods string   `json:"fundraisingMethods" gorm:"fundraising_methods" bson:"fundraising_methods" title:"募集方式"`
	ProductStatus      string   `json:"productStatus" gorm:"product_status" bson:"product_status" title:"产品状态"`
	LicenseNumber      string   `json:"licenseNumber" gorm:"license_number" bson:"license_number" title:"备案号"`
	ProductDesc        string   `json:"productDesc" gorm:"product_desc" bson:"product_desc" title:"产品描述"`
	Tags               []string `json:"tags" gorm:"type:text" bson:"tags" title:"标签"`
	Remark             string   `json:"remark" gorm:"remark" bson:"remark" title:"备注"`
}

func NewProduct() *Product {
	return &Product{}
}