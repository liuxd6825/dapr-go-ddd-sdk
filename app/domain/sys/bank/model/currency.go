package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Bank struct {
	xbase.BaseModel `bson:",inline"`
	Name            string `json:"name" gorm:"name" bson:"name"  index:"asc"  title:"币种名称"`
	Keywords        string `json:"keywords" gorm:"keywords" bson:"keywords"  title:"关键字"`
	Order           int    `json:"order" gorm:"order" bson:"order" index:"asc" title:"顺序"`
}

func NewBank() (*Bank, error) {
	return &Bank{}, nil
}

func (cur *Bank) GetName() string {
	return cur.Name
}

func (cur *Bank) GetKeywords() string {
	return cur.Keywords
}
