package schema

import (
	"context"
	"reflect"
)

type Field struct {
	Name           string   // 属性名
	DBName         string   // 数据库字段名
	DataType       DataType // 数据类型
	PrimaryKey     bool     // 是主键
	NotNull        bool     // 不能为空
	Unique         bool     // 唯一键值
	Comment        string   // 字段说明
	Size           int      // 字段大小
	Creatable      bool     // 是否可创建
	Updatable      bool     // 是否可更新
	Readable       bool     // 可读取的字段
	AutoCreateTime TimeType
	AutoUpdateTime TimeType

	BindNames              []string
	EmbeddedBindNames      []string
	GORMDataType           DataType
	AutoIncrement          bool
	AutoIncrementIncrement int64
	HasDefaultValue        bool
	DefaultValue           string
	DefaultValueInterface  interface{}
	Precision              int
	Scale                  int
	TagSettings            map[string]string
	ReflectValueOf         func(context.Context, reflect.Value) reflect.Value
	ValueOf                func(context.Context, reflect.Value) (value interface{}, zero bool) // 转换实体属性到数据库字段
	Set                    func(context.Context, reflect.Value, interface{}) error             // 转换数据库值到实体属性类型
}

func NewField() *Field {
	return &Field{
		Creatable: true,
		Updatable: true,
		Readable:  true,
	}
}

func (field *Field) SetName(val string) *Field {
	field.Name = val
	return field
}

func (field *Field) SetDbName(val string) *Field {
	field.DBName = val
	return field
}

func (field *Field) SetDataType(val DataType) *Field {
	field.DataType = val
	return field
}

func (field *Field) SetSize(val int) *Field {
	field.Size = val
	return field
}
