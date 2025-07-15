package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type Draw struct {
	xbase.BaseModel `bson:",inline"`
	Status          string `json:"status" gorm:"status" bson:"status"`
	Name            string `json:"name" gorm:"name" bson:"name"`
	FileName        string `json:"fileName" gorm:"file_name" bson:"file_name"`
}
