package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type App struct {
	xbase.BaseModel `bson:",inline"`
	Name            string `json:"name"  gorm:"name" bson:"name"`
	Code            string `json:"code" gorm:"code"  bson:"code"`
}

func NewApp() (*App, error) {
	return &App{}, nil
}

type AppTree struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Children []*AppTree `json:"children"`
}
