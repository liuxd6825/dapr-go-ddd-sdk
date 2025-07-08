package model

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/enums"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/xbase/xmodel"
)

type Template struct {
	xmodel.Base `gorm:",inline"`
	Name        string           `json:"name" gorm:"name"  gorm:"name" index:"" `
	FileId      string           `json:"fileId" gorm:"file_id"   gorm:"file_id" index:"" `
	FileName    string           `json:"fileName"gorm:"file_name"  gorm:"file_name" index:"" `
	SheetName   string           `json:"sheetName"  gorm:"sheet_name" gorm:"sheet_name"`
	BankName    string           `json:"bankName" gorm:"bank_name"  gorm:"bank_name" index:"" `
	MasterType  enums.MasterType `json:"masterType" gorm:"master_type"  gorm:"master_type"  index:"" `
	MapHeads    []*MapHead       `json:"mapHeads" gorm:"map_heads" gorm:"map_heads"`
	Fields      []*TemplateField `json:"fields" gorm:"fields" gorm:"fields"`
}

func NewTemplate() (*Template, error) {
	return &Template{}, nil
}

type MapHead struct {
	HeadRow int64            `json:"headRow" gorm:"head_row"   gorm:"head_row"`
	Columns []*MapHeadColumn `json:"columns"  gorm:"columns"  gorm:"columns"`
}

// MapHeadColumn
// @Description: 映射规则
type MapHeadColumn struct {
	Key      string   `json:"key"  gorm:"key"  gorm:"key"`
	Label    string   `json:"label"  gorm:"label"  gorm:"label"`
	DataType DataType `json:"dataType"   gorm:"data_type"`
}
