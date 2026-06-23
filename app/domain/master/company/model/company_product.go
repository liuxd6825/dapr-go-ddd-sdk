package model

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"
)

// CompanyProduct 公司-产品关联
type CompanyProduct struct {
	xbase.BaseModel `bson:",inline"`

	CompanyId    string     `json:"companyId" gorm:"company_id" bson:"company_id" index:"" title:"公司ID"`
	RelationType string     `json:"relationType" gorm:"relation_type" bson:"relation_type" title:"关系类型"`
	Name         string     `json:"name" gorm:"name" bson:"name" title:"产品名称"`
	Code         string     `json:"code" gorm:"code" bson:"code" title:"产品编号"`
	Nature       string     `json:"nature" gorm:"nature" bson:"nature" title:"产品性质"`
	ProductType  string     `json:"productType" gorm:"product_type" bson:"product_type" title:"产品类型"`
	Seller       string     `json:"seller" gorm:"seller" bson:"seller" title:"发售方"`
	StartDate    *time.Time `json:"startDate" gorm:"start_date" bson:"start_date" title:"生效日期"`
	EndDate      *time.Time `json:"endDate" gorm:"end_date" bson:"end_date" title:"失效日期"`
}

func NewCompanyProduct() *CompanyProduct {
	return &CompanyProduct{}
}
