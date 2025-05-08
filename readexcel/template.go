package readexcel

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Template
// @Description: 导入模板
type Template struct {
	Id         string     `json:"id"`
	Consts     []*Const   `json:"consts"`
	Fields     []*Field   `json:"fields"`
	Heads      []*MapHead `json:"columns"`
	Script     *string    `json:"script"`
	constsMap  map[string]*Const
	fieldsMap  map[string]*Field
	columnsMap map[string]*MapColumn
}

type TemplateBuilder interface {
	NewTemplate()
}

// NewTemplate
// @Description: 新建
// @return *Template
func NewTemplate(fields []*Field, heads []*MapHead, consts []*Const) (*Template, error) {
	temp := &Template{
		Fields:     fields,
		Heads:      heads,
		Consts:     consts,
		constsMap:  make(map[string]*Const),
		fieldsMap:  make(map[string]*Field),
		columnsMap: make(map[string]*MapColumn),
	}
	if err := temp.Init(); err != nil {
		return nil, err
	}
	return temp, nil
}

// MapHead
// @Description: MapHead
type MapHead struct {
	RowNum  int64        `json:"rowNum"`
	Columns []*MapColumn `json:"columns"`
}

// DataType
// @Description: 数据类型
type DataType string

const (
	None     DataType = ""
	String   DataType = "str"
	Integer  DataType = "int"
	DateTime DataType = "dateTime"
	Date     DataType = "date"
	Time     DataType = "time"
	Money    DataType = "money"
)

// Field
// @Description: 转成字段项
type Field struct {
	Name      string     `json:"name,omitempty"`
	Title     string     `json:"title,omitempty"`
	DataType  DataType   `json:"dataType,omitempty"`
	Size      int64      `json:"size,omitempty"`
	AllowNull bool       `json:"allowNull,omitempty"`
	Validator string     `json:"validator,omitempty"`
	MapKeys   []string   `json:"mapKeys,omitempty"`
	DefFloat  *float64   `json:"defFloat,omitempty"`
	DefInt    *int64     `json:"defInt,omitempty"`
	DefString string     `json:"defString,omitempty"`
	DefTime   *time.Time `json:"defTime,omitempty"`
	Script    string     `json:"script,omitempty"`
	Replaces  []Replace  `json:"replaces,omitempty"`
}

type Replace struct {
	Old string `json:"old"`
	New string `json:"new"`
}

func NewField(name string, title string, dataType DataType, allowNull bool, keys ...string) *Field {
	return &Field{
		Name:      strings.ReplaceAll(name, "\t", ""),
		Title:     strings.ReplaceAll(title, "\t", ""),
		DataType:  dataType,
		Size:      0,
		AllowNull: allowNull,
		Validator: "",
		MapKeys:   keys,
	}
}

// MapItem
// @Description:
type MapItem struct {
	MapType MapType `json:"mapType"`
	Name    string  `json:"name"`
}

type ConstType int64

const (
	ConstTypeCell ConstType = iota
	ConstTypeValue
)

// Const
// @Description: 常量
type Const struct {
	Key   string    `json:"key"`
	Type  ConstType `json:"type"`
	Point *Point    `json:"point"`
	Value *string   `json:"value"`
}

// Point
// @Description: 位置
type Point struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// MapType
// @Description: 映射模式
type MapType int64

const (
	MapTypeColumn MapType = iota // 单元格模式
	MapTypeConst                 // 常量
)

// MapColumn
// @Description: 映射规则
type MapColumn struct {
	Key      string   `json:"key"`
	Label    string   `json:"colLabel"`
	DataType DataType `json:"dataType"`
	Script   string   `json:"script"`
}

type KeyValueOptions struct {
	Temp *Template
	//Sheet    *xlsx.Sheet
	DataRow  *DataRow
	RowCells map[string]*DataCell
}

type KeyValue interface {
	GetKey() string
	GetValue(c *KeyValueOptions) string
}

// NotFoundKeyError
// @Description:
type NotFoundKeyError struct {
	typeName string
	key      string
}

func newNotFoundKeyError(typeName string, key string) error {
	return &NotFoundKeyError{typeName: typeName, key: key}
}

// Validator
//
//	@Description: 验证器
type Validator func(cellText string) (any, error)
type SetFieldFunc func(keys []string, values []string) (string, error)

var validatorRegistry = make(map[string]Validator)
var setFieldRegistry = make(map[string]SetFieldFunc)

func setFieldJoin(keys []string, values []string) (string, error) {
	return strings.Join(values, ""), nil
}

func GetFieldFunc(key string) (SetFieldFunc, bool) {
	if f, ok := setFieldRegistry[key]; ok {
		return f, true
	}
	return nil, false
}

func AddFieldFunc(key string, fun SetFieldFunc) error {
	setFieldRegistry[key] = fun
	return nil
}

// AddValidator
// @Description: 添加验证器
// @param key  关键字
// @param validator 验证器
func AddValidator(key string, validator Validator) {
	validatorRegistry[key] = validator
}

func (t *Template) Init() error {
	t.fieldsMap = make(map[string]*Field)
	t.constsMap = make(map[string]*Const)
	t.columnsMap = make(map[string]*MapColumn)
	for _, field := range t.Fields {
		if field != nil {
			t.fieldsMap[field.Name] = field
		}

	}
	for _, con := range t.Consts {
		if con != nil {
			t.constsMap[con.Key] = con
		}
	}
	for _, row := range t.Heads {
		if len(row.Columns) > 0 {
			for _, col := range row.Columns {
				if col != nil {
					t.columnsMap[getCellString(col.Key)] = col
				}
			}
		}
	}
	return nil
}

// AddConst
// @Description: 添加变量
// @receiver t
// @param name
// @param rowNum
// @param colNum
// @return error
func (t *Template) AddConst(name string, ctype ConstType, value string, point *Point) error {
	c := &Const{Key: name, Type: ctype, Point: point, Value: PString(value)}
	if _, ok := t.constsMap[name]; !ok {
		t.Consts = append(t.Consts, c)
		t.constsMap[c.Key] = c
		return nil
	}
	return errors.New(name)
}

func (t *Template) GetMapColumn(key string) (*MapColumn, error) {
	if c, ok := t.columnsMap[key]; ok {
		return c, nil
	}
	return nil, newNotFoundKeyError("MapColumn", key)
}

func (t *Template) GetField(key string) *Field {
	if c, ok := t.fieldsMap[key]; ok {
		return c
	}
	return nil
}

func (t *Template) GetKeyValue(key string) KeyValue {
	if c, ok := t.constsMap[key]; ok {
		return c
	} else if c, ok := t.columnsMap[key]; ok {
		return c
	}
	return nil
}
func (t *Template) GetConst(key string) (*Const, error) {
	if c, ok := t.constsMap[key]; ok {
		return c, nil
	}
	return nil, newNotFoundKeyError("Const", key)
}

func (t *Template) ValueToFloat(fieldKey string, row *DataRow, defValue *float64) (*float64, error) {
	return ValueToFloat(row.Values[fieldKey], defValue)
}

func (t *Template) ValueToDate(fieldKey string, row *DataRow, defValue *time.Time) (*time.Time, error) {
	return ValueToDate(row.Values[fieldKey], defValue)
}

func (t *Template) ValueToInt(fieldKey string, row *DataRow, defValue *int64) (*int64, error) {
	return ValueToInt(row.Values[fieldKey], defValue)
}

func (t *Template) ValueToBool(fieldKey string, row *DataRow, defValue bool) (bool, error) {
	return ValueToBool(row.Values[fieldKey], defValue)
}

func (t *Template) ValueToString(fieldKey string, row *DataRow) (res string, err error) {
	f, _ := t.fieldsMap[fieldKey]
	return ValueToString(row.Values[fieldKey], "", f.GetReplace()...)
}

func (c *MapColumn) GetKey() string {
	return c.Key
}

func (c *MapColumn) GetValue(ctx *KeyValueOptions) string {
	cell := ctx.RowCells[c.Key]
	if cell != nil {
		return cell.Value
	}
	return ""
}

func (c *MapColumn) GetColNum() int64 {
	return LabelToInt(c.Label)
}

func (c *Const) GetKey() string {
	return c.Key
}

func (c *Const) GetValue(ctx *KeyValueOptions) string {
	if c.Type == ConstTypeCell {
		if c.Value != nil {
			return *c.Value
		}
		/*
			cell := ctx.Sheet.Cell(c.Point.Row, c.Point.Col)
			if cell != nil {
				c.Value = &cell.Value
			} else {
				c.Value = PString("")
			} */
		return *c.Value
	} else if c.Type == ConstTypeValue {
		return *c.Value
	}
	return ""
}

func (n *NotFoundKeyError) Error() string {
	return fmt.Sprintf("not found key '%s'", n.key)
}

func (f *Field) SetName(v string) *Field {
	f.Name = v
	return f
}
func (f *Field) SetTitle(v string) *Field {
	f.Title = v
	return f
}
func (f *Field) SetDataType(v DataType) *Field {
	f.DataType = v
	return f
}
func (f *Field) SetSize(v int64) *Field {
	f.Size = v
	return f
}
func (f *Field) SetAllowNull(v bool) *Field {
	f.AllowNull = v
	return f
}
func (f *Field) SetValidator(v string) *Field {
	f.Validator = v
	return f
}
func (f *Field) SetMapKeys(v ...string) *Field {
	var list []string
	for _, val := range v {
		list = append(list, getCellString(val))
	}
	f.MapKeys = list
	return f
}

func (f *Field) SetDefFloat(v *float64) *Field {
	f.DefFloat = v
	return f
}

func (f *Field) GetDefString() string {
	if f == nil {
		return ""
	}
	return f.DefString
}

func (f *Field) GetDefFloat() *float64 {
	return f.DefFloat
}

func (f *Field) GetDefInt() *int64 {
	return f.DefInt
}

func (f *Field) GetDefTime() *time.Time {
	return f.DefTime
}

func (f *Field) SetScript(v string) *Field {
	f.Script = v
	return f
}
func (f *Field) GetReplace() []Replace {
	if f == nil {
		return []Replace{}
	}
	return f.Replaces
}

func (f *Field) AddReplace(old string, new string) *Field {
	f.Replaces = append(f.Replaces, Replace{Old: old, New: new})
	return f
}

func (f *Field) AddReplaces(s ...Replace) *Field {
	f.Replaces = append(f.Replaces, s...)
	return f
}

func DefaultValue(v string) *string {
	return &v
}

func PString(v string) *string {
	return &v
}

func LabelToInt(label string) int64 {
	// "FA=157=6*26+1=156+1", "FR=173"
	var count int32 = 0
	lab := strings.ToLower(label)
	l := int32(len(lab))
	for i, k := range lab {
		p := l - int32(i+1)
		v := k - 96
		if p > 0 {
			v = v * 26
		}
		count += v
	}
	return int64(count - 1)
}

func GetCellLabel(idx int) string {
	res := ""
	v := idx / 26
	m := idx % 26

	if idx <= 26 {
		res = res + label(m)
	} else if v > 0 && m == 0 {
		res = label(v-1) + "Z"
	} else {
		res = res + label(v)
		res = res + label(m)
	}
	return res
}

func label(v int) string {
	v = v % 26
	if v > 0 {
		v = v + 64
		b := byte(v)
		return fmt.Sprintf("%c", b)
	}
	return "Z"
}
