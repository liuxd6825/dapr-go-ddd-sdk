package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/xcommon/xbase"
)

type Template struct {
	xbase.BaseModel `bson:",inline"`
	TaskId          string           `json:"taskId" gorm:"task_id"  bson:"task_id" index:""`
	Name            string           `json:"name" gorm:"name"  bson:"name" index:"" `
	FileId          string           `json:"fileId" gorm:"file_id"   bson:"file_id" index:"" `
	FileName        string           `json:"fileName" gorm:"file_name"  bson:"file_name" index:"" `
	SheetName       string           `json:"sheetName"  gorm:"sheet_name" bson:"sheet_name"`
	SheetId         string           `json:"sheetId"  gorm:"sheet_id" bson:"sheet_id"`
	BankName        string           `json:"bankName" gorm:"bank_name"  bson:"bank_name"`
	SchemaId        string           `json:"schemaId" gorm:"schema_id"  bson:"schema_id"  index:"" `
	SchemaName      string           `json:"schemaName" gorm:"schema_name"  bson:"schema_name"  `
	MapHeads        []*MapHead       `json:"mapHeads" gorm:"map_heads;type:json" bson:"map_heads"`
	Fields          []*TemplateField `json:"fields" gorm:"fields;type:json" bson:"fields"`
}

func NewTemplate() (*Template, error) {
	return &Template{}, nil
}

type MapHead struct {
	HeadRow int64            `json:"headRow" gorm:"head_row"   bson:"head_row"`
	Columns []*MapHeadColumn `json:"columns"  gorm:"columns"  bson:"columns"`
}

// MapHeadColumn
// @Description: 映射规则
type MapHeadColumn struct {
	Key      string   `json:"key"  gorm:"key"  bson:"key"`
	Label    string   `json:"label"  gorm:"label"  bson:"label"`
	DataType DataType `json:"dataType"   gorm:"data_type" bson:"data_type"`
}

// TemplateField
// @Description: 转成字段项
type TemplateField struct {
	Key       string   `json:"key" gorm:"key" bson:"key"`
	Name      string   `json:"name" gorm:"name" bson:"name"`
	MapKeys   []string `json:"mapKeys" gorm:"map_keys" bson:"map_keys"`
	Script    string   `json:"script"  gorm:"script" bson:"script"`
	AllowNull *bool    `json:"allowNull" gorm:"allow_null" bson:"allow_null"`
	IsHide    *bool    `json:"isHide" gorm:"is_hide" bson:"is_hide"`
}

// DataType
// @Description: 数据类型
type DataType = readexcel.DataType
