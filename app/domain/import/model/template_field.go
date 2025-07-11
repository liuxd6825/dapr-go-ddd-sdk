package model

// TemplateField
// @Description: 转成字段项
type TemplateField struct {
	Key       string   `json:"key,omitempty" gorm:"key" bson:"key"`
	Name      string   `json:"name,omitempty" gorm:"name" bson:"name"`
	MapKeys   []string `json:"mapKeys,omitempty" gorm:"map_keys" bson:"map_keys"`
	Script    string   `json:"script,omitempty"  gorm:"script" bson:"script"`
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
