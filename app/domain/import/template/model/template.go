package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/recordie/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase"
)

type Template struct {
	xbase.BaseModel `gorm:",inline"`
	Name            string           `json:"name" gorm:"name"  bson:"name" index:"" `
	FileId          string           `json:"fileId" gorm:"file_id"   bson:"file_id" index:"" `
	FileName        string           `json:"fileName"gorm:"file_name"  bson:"file_name" index:"" `
	SheetName       string           `json:"sheetName"  gorm:"sheet_name" bson:"sheet_name"`
	BankName        string           `json:"bankName" gorm:"bank_name"  bson:"bank_name" index:"" `
	MasterType      enums.MasterType `json:"masterType" gorm:"master_type"  bson:"master_type"  index:"" `
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
