package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Fun struct {
	xbase.BaseModel `bson:",inline"`
	AppId           string `json:"appId" gorm:"app_id" bson:"app_id"`
	ParentId        string `json:"parentId" gorm:"parent_id" bson:"parent_id"`
	Name            string `json:"name" gorm:"name" bson:"name"`
	Code            string `json:"code" gorm:"code"  bson:"code"`
	OrderNum        int64  `json:"orderNum" gorm:"order_num"  bson:"order_num"`
}

func NewFun() (*Fun, error) {
	return &Fun{}, nil
}

type FunView struct {
	xbase.BaseModel `bson:",inline"`
	AppId           string      `json:"appId" gorm:"app_id" bson:"app_id"`
	Name            string      `json:"name" gorm:"name" bson:"name"`
	Path            []string    `json:"path" gorm:"path" bson:"path"`
	Cards           []*CardFile `json:"cards" gorm:"cards" bson:"cards"`
	OrderNum        float64     `json:"orderNum" gorm:"order_num"  bson:"order_num"`
}
