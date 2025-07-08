package model

// TemplateField
// @Description: 转成字段项
type TemplateField struct {
	Key       string   `json:"key,omitempty" gorm:"key" gorm:"key"`
	Name      string   `json:"name,omitempty" gorm:"name" gorm:"name"`
	MapKeys   []string `json:"mapKeys,omitempty" gorm:"map_keys" gorm:"map_keys"`
	Script    string   `json:"script,omitempty"  gorm:"script" gorm:"script"`
	AllowNull *bool    `json:"allowNull" gorm:"allow_null" gorm:"allow_null"`
	IsHide    *bool    `json:"isHide" gorm:"is_hide" gorm:"is_hide"`
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
