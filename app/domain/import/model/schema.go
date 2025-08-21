package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type SchemaModel struct {
	xbase.BaseModel
	Content string `json:"content" gorm:"content" bson:"content" title:"内容"`
}
