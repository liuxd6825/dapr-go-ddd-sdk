package ddd_sql

import (
	"gorm.io/gorm"
	"strings"
)

// ConvertCamelToSnake 将驼峰命名转换为下划线命名（snake_case）
func ConvertCamelToSnake(str string) string {
	var result []rune
	for i, s := range str {
		if i > 0 && (s >= 'A' && s <= 'Z') {
			result = append(result, '_')
		}
		result = append(result, s)
	}
	return strings.ToLower(string(result))
}

// ConvertSnakeToCamel 将下划线命名转换为驼峰命名（camelCase）
func ConvertSnakeToCamel(str string) string {
	parts := strings.Split(str, "_")
	for i, part := range parts {
		if i > 0 {
			parts[i] = strings.Title(part)
		}
	}
	return strings.Join(parts, "")
}

// GormFieldPlugin 插件结构体
type GormFieldPlugin struct{}

// NewGormFieldPlugin 插件初始化函数
func NewGormFieldPlugin() *GormFieldPlugin {
	return &GormFieldPlugin{}
}

func (f *GormFieldPlugin) Name() string {
	return "GormFieldPlugin"
}

func (f *GormFieldPlugin) Initialize(db *gorm.DB) error {
	return nil
}

// AfterFind 在查询结果中转换字段名
func (f *GormFieldPlugin) AfterFind(tx *gorm.DB) (err error) {
	// 遍历数据库字段名，将其转为驼峰命名
	for i := range tx.Statement.Schema.Fields {
		tx.Statement.Schema.Fields[i].DBName = ConvertCamelToSnake(tx.Statement.Schema.Fields[i].Name)
	}
	return nil
}

// BeforeCreate 在插入数据时转换表字段的值 (字段插入)
func (f *GormFieldPlugin) BeforeCreate(tx *gorm.DB) (err error) {
	// 遍历并转换数据键名称为下划线命名方式
	for key := range tx.Statement.Clauses {
		tx.Statement.Clauses[ConvertCamelToSnake(key)] = tx.Statement.Clauses[key]
	}
	return nil
}
