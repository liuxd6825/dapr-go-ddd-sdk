package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Group struct {
	xbase.BaseModel `bson:",inline"`
	HomeId          string `json:"homeId" gorm:"home_id" bson:"home_id"`
	Name            string `json:"name" gorm:"name" bson:"name"`
	OrderNum        int64  `json:"orderNum" gorm:"order_num"  bson:"order_num"`
	SourceId        string `json:"sourceId" gorm:"source_id"  bson:"source_id"`
	IsHide          bool   `json:"isHide" gorm:"is_hide" bson:"is_hide"`
}

func NewGroup() (*Group, error) {
	return &Group{}, nil
}
