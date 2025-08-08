package rsql

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

// Builder RSQL语包生成器
type Builder struct {
	conditions []Condition
}

type NewBuilderOptions func(b *Builder) Condition

func NewBuilder(opts ...NewBuilderOptions) *Builder {
	b := &Builder{}
	var conditions []Condition
	for _, opt := range opts {
		c := opt(b)
		if c != nil {
			conditions = append(conditions, c)
		}
		b.conditions = conditions
	}
	return b
}

func (b *Builder) RSQL(filter string) Condition {
	return &baseCondition{
		field:    "",
		operator: "",
		value:    filter,
	}
}

// formatValue 处理不同数据类型的值格式化
func (b *Builder) formatValue(value interface{}) string {
	return formatValue(value)
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

func (b *Builder) Like(field string, value interface{}) Condition {
	return Like(field, value)
}

// Like 模糊
func Like(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: "~=",
		value:    formatValue(value),
	}
}

func (b *Builder) Eq(field string, value interface{}) Condition {
	return Eq(field, value)
}

// Eq 等于
func Eq(field string, value interface{}) Condition {
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

func (b *Builder) Null(field string) Condition {
	return Null(field)
}

// Null 等于
func Null(field string) Condition {
	return &baseCondition{
		field:    field,
		operator: "=null=",
		value:    "0",
	}
}

func (b *Builder) NotNull(field string) Condition {
	return NotNull(field)
}

// NotNull 等于
func NotNull(field string) Condition {
	return &baseCondition{
		field:    field,
		operator: "=!null=",
		value:    "0",
	}
}

func (b *Builder) Start(field string, val any) Condition {
	return Start(field, val)
}

// Start 以..开头
func Start(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=start=",
		value:    formatValue(val),
	}
}

// End 以..结尾
func (b *Builder) End(field string, val any) Condition {
	return End(field, val)
}

func End(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=end=",
		value:    formatValue(val),
	}
}

func (b *Builder) Contains(field string, val any) Condition {
	return Contains(field, val)
}

// Contains 包含
func Contains(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=contains=",
		value:    formatValue(val),
	}
}

func (b *Builder) NotContains(field string, val any) Condition {
	return NotContains(field, val)
}

// NotContains 包含
func NotContains(field string, val any) Condition {
	return &baseCondition{
		field:    field,
		operator: "=!contains=",
		value:    formatValue(val),
	}
}

// Neq 不等于
func (b *Builder) Neq(field string, val interface{}) Condition {
	return Neq(field, val)
}

func Neq(field string, val any) Condition {
	if val == nil {
		return &baseCondition{
			field:    field,
			operator: "=!null=",
			value:    "0",
		}
	}
	return &baseCondition{
		field:    field,
		operator: "!=",
		value:    formatValue(val),
	}
}

func (b *Builder) Gt(field string, value interface{}) Condition {
	return Gt(field, value)
}

// Gt 大于
func Gt(field string, value interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: ">",
		value:    formatValue(value),
	}
}

func (b *Builder) Ge(field string, val interface{}) Condition {
	return Ge(field, val)
}

// Ge 大于等于
func Ge(field string, val interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: ">=",
		value:    formatValue(val),
	}
}

func (b *Builder) Lt(field string, val interface{}) Condition {
	return Lt(field, val)
}

// Lt 小于
func Lt(field string, val interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: "<",
		value:    formatValue(val),
	}
}

func (b *Builder) Le(field string, val interface{}) Condition {
	return Lt(field, val)
}

// Le 小于等于
func Le(field string, val interface{}) Condition {
	return &baseCondition{
		field:    field,
		operator: "<=",
		value:    formatValue(val),
	}
}

func (b *Builder) In(field string, val interface{}) Condition {
	return In(field, val)
}

// In 条件生成方法
func In(field string, values interface{}) Condition {
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

func (b *Builder) Out(field string, value interface{}) Condition {
	return Out(field, value)
}

// Out 条件生成方法
func Out(field string, values interface{}) Condition {
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

// And 与逻辑组合方法
func (b *Builder) And(conditions ...Condition) Condition {
	return And(conditions...)
}
func And(conditions ...Condition) Condition {
	return &compositeCondition{
		operator:   " and ",
		conditions: conditions,
	}
}

// Or 或逻辑
func (b *Builder) Or(conditions ...Condition) Condition {
	return Or(conditions...)
}

func Or(conditions ...Condition) Condition {
	return &compositeCondition{
		operator:   " or ",
		conditions: conditions,
	}
}

// GetList 取得slice中元素指定字段的值，支持struct与map类型。
func (b *Builder) GetList(list any, filed string) []string {
	res, err := GetFieldValues(list, filed)
	if err != nil {
		panic(err)
	}
	return res
}

func (b *Builder) Build() string {
	sb := strings.Builder{}
	for _, item := range b.conditions {
		sb.WriteString(item.Build())
	}
	return sb.String()
}
