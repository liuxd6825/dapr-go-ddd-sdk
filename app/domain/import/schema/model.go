package schema

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"

type SchemaModel struct {
	xbase.BaseModel
	Content string `json:"content" gorm:"content" bson:"content" title:"内容"`
}
