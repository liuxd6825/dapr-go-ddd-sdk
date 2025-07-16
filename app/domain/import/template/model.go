package template

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type Template struct {
	xbase.BaseModel `gorm:",inline"`
	Name            string           `json:"name" gorm:"name"  bson:"name" index:"" `
	FileId          string           `json:"fileId" gorm:"file_id"   bson:"file_id" index:"" `
	FileName        string           `json:"fileName"gorm:"file_name"  bson:"file_name" index:"" `
	SheetName       string           `json:"sheetName"  gorm:"sheet_name" bson:"sheet_name"`
	BankName        string           `json:"bankName" gorm:"bank_name"  bson:"bank_name" index:"" `
	SchemaId        string           `json:"schemaId" gorm:"schema_id"  bson:"schema_id"  index:"" `
	MapHeads        []*MapHead       `json:"mapHeads" gorm:"map_heads" bson:"map_heads"`
	Fields          []*TemplateField `json:"fields" gorm:"fields" bson:"fields"`
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
type DataType string

const (
	DataType_None     DataType = ""
	DataType_String   DataType = "str"
	DataType_Integer  DataType = "int"
	DataType_DateTime DataType = "dateTime"
	DataType_Date     DataType = "date"
	DataType_Time     DataType = "time"
	DataType_Money    DataType = "money"
)
