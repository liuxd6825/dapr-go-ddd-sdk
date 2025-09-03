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
}

func NewFun() (*Fun, error) {
	return &Fun{}, nil
}
