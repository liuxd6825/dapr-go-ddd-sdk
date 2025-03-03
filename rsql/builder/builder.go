package builder

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Condition 接口定义查询条件
type Condition interface {
	Build() string
}

// 基础比较条件
type baseCondition struct {
	field    string
	operator string
	value    string
}

func (c *baseCondition) Build() string {
	return fmt.Sprintf("%s%s%s", c.field, c.operator, c.value)
}

// 组合条件
type compositeCondition struct {
	operator   string
	conditions []Condition
}

func (c *compositeCondition) Build() string {
	var parts []string
	for _, cond := range c.conditions {
		parts = append(parts, cond.Build())
	}
	return fmt.Sprintf("(%s)", strings.Join(parts, c.operator))
}

// IN 条件
type inCondition struct {
	field  string
	values string
}

func (c *inCondition) Build() string {
	return fmt.Sprintf("%s=in=%s", c.field, c.values)
}

type outCondition struct {
	field  string
	values string
}

func (c *outCondition) Build() string {
	return fmt.Sprintf("%s=out=%s", c.field, c.values)
}

// Builder 主结构体
type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

// formatValue 处理不同数据类型的值格式化
func formatValue(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case string:
		escaped := strings.ReplaceAll(v, `"`, `""`)
		return fmt.Sprintf(`'%s'`, escaped)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%v", v)
	case time.Time:
		return fmt.Sprintf(`"%s"`, v.Format(time.RFC3339))
	default:
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Ptr {
			if val.IsNil() {
				return "null"
			}
			return formatValue(val.Elem().Interface())
		}
		return fmt.Sprintf(`"%v"`, v)
	}
}

// Like 模糊
func (b *Builder) Like(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: "~=",
		value:    formatValue(value),
	}
}

// Eq 等于
func (b *Builder) Eq(field string, value interface{}) Condition {
	if value == nil {
		return &baseCondition{
			field:    field,
			operator: "=null=",
			value:    "0",
		}
	}
	return &baseCondition{
		field:    field,
		operator: "==",
		value:    formatValue(value),
	}
}

// Null 等于
func (b *Builder) Null(field string) Condition {
	return &baseCondition{
		field:    field,
		operator: "=null=",
		value:    "0",
	}
}

// NotNull 等于
func (b *Builder) NotNull(field string) Condition {
	return &baseCondition{
		field:    field,
		operator: "=!null=",
		value:    "0",
	}
}

// Start 以..开头
func (b *Builder) Start(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=start=",
		value:    formatValue(val),
	}
}

// End 以..结尾
func (b *Builder) End(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=end=",
		value:    formatValue(val),
	}
}

// Contains 包含
func (b *Builder) Contains(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=contains=",
		value:    formatValue(val),
	}
}

// NotContains 包含
func (b *Builder) NotContains(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=!contains=",
		value:    formatValue(val),
	}
}

// Neq 不等于
func (b *Builder) Neq(field string, value interface{}) Condition {
	if value == nil {
		return &baseCondition{
			field:    field,
			operator: "=!null=",
			value:    "0",
		}
	}
	return &baseCondition{
		field:    field,
		operator: "!=",
		value:    formatValue(value),
	}
}

// Gt 大于
func (b *Builder) Gt(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: ">",
		value:    formatValue(value),
	}
}

// Ge 大于等于
func (b *Builder) Ge(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: ">=",
		value:    formatValue(value),
	}
}

func (b *Builder) Lt(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: "<",
		value:    formatValue(value),
	}
}

func (b *Builder) Le(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: "<=",
		value:    formatValue(value),
	}
}

// In 条件生成方法
func (b *Builder) In(field string, values interface{}) Condition {
	val := reflect.ValueOf(values)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		panic("In values must be a slice or array")
	}

	var formatted []string
	for i := 0; i < val.Len(); i++ {
		formatted = append(formatted, formatValue(val.Index(i).Interface()))
	}

	return &inCondition{
		field:  field,
		values: fmt.Sprintf("(%s)", strings.Join(formatted, ",")),
	}
}

// Out 条件生成方法
func (b *Builder) Out(field string, values interface{}) Condition {
	val := reflect.ValueOf(values)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		panic("In values must be a slice or array")
	}

	var formatted []string
	for i := 0; i < val.Len(); i++ {
		formatted = append(formatted, formatValue(val.Index(i).Interface()))
	}

	return &outCondition{
		field:  field,
		values: fmt.Sprintf("(%s)", strings.Join(formatted, ",")),
	}
}

// 逻辑组合方法
func (b *Builder) And(conditions ...Condition) Condition {
	return &compositeCondition{
		operator:   " and ",
		conditions: conditions,
	}
}

func (b *Builder) Or(conditions ...Condition) Condition {
	return &compositeCondition{
		operator:   " or ",
		conditions: conditions,
	}
}
