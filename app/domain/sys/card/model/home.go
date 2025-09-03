package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/sys/card/enum"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Home struct {
	xbase.BaseModel `bson:",inline"`
	Name            string        `json:"name"  gorm:"name" bson:"name"`
	Code            string        `json:"code" gorm:"code"  bson:"code"`
	HomeType        enum.HomeType `json:"homeType" gorm:"home_type"  bson:"home_type"`
	IsDefault       bool          `json:"isDefault" gorm:"is_default" bson:"is_default"`
}

func NewHome() (*Home, error) {
	return &Home{}, nil
}

type HomeView struct {
	Home   *Home    `json:"home"`
	Groups []*Group `json:"groups"`
	Cards  []*Card  `json:"cards"`
}

func NewHomeView() *HomeView {
	return &HomeView{}
}
