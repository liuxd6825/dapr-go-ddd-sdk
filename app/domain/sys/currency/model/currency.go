package model

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Currency struct {
	xbase.BaseModel `bson:",inline"`
	Area            string  `json:"area" gorm:"area" bson:"area"  index:"asc"  title:"国家或地区" `
	Name            string  `json:"name" gorm:"name" bson:"name"  index:"asc"  title:"币种名称"`
	Code1           string  `json:"code1" gorm:"code1" bson:"code1"  index:"asc"  title:"编码1"`
	Code2           string  `json:"code2" gorm:"code2" bson:"code2"  index:"asc"  title:"编码2"`
	Code3           string  `json:"code3" gorm:"code3" bson:"code3"  index:"asc"  title:"编码3"`
	Rate            float64 `json:"rate" gorm:"rate" bson:"rate" title:"对美元汇率"`
	CnRate          float64 `json:"cnRate" gorm:"cn_rate" bson:"cn_rate" title:"对人民币汇率"`
	Order           int     `json:"order" gorm:"order" bson:"order" index:"asc" title:"顺序"`
}

func NewCurrency() (*Currency, error) {
	return &Currency{}, nil
}

func (cur *Currency) GetName() string {
	return cur.Name
}

func (cur *Currency) GetKeywords() string {
	return fmt.Sprintf("%s,%s,%s,", cur.Code1, cur.Code2, cur.Code3)
}
