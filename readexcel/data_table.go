package readexcel

import (
	"errors"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/setting"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DataTable
// @Description: 数据表格
type DataTable struct {
	Temp   *Template  `json:"temp"`
	Rows   []*DataRow `json:"rows"`
	Errors []int64    `json:"errors"`
}

// DataRow
// @Description: 数据行
type DataRow struct {
	RowNum int64                `json:"rowNum"` // 行号
	Values map[string]any       `json:"values"` // 数据值
	Cells  map[string]DataCells `json:"cells"`  // 单元格值
	Errors map[string][]string  `json:"errors"` // 错误
}

type DataCells []*DataCell

type DataCell struct {
	Key    string   `json:"key" bson:"key,omitempty"`
	Value  string   `json:"value" bson:"value,omitempty"`
	Col    int64    `json:"col" bson:"col,omitempty"`
	Row    int64    `json:"row" bson:"row,omitempty"`
	Errors []string `json:"errors" bson:"errors,omitempty"`
}

const BaseFormat = "2006-01-02 15:04:05"

func NewDataTable(temp *Template) *DataTable {
	return &DataTable{
		Temp: temp,
		Rows: make([]*DataRow, 0),
	}
}

func NewDataRow() *DataRow {
	return &DataRow{
		Cells:  make(map[string]DataCells),
		Values: make(map[string]any),
		Errors: map[string][]string{},
	}
}

func (t *DataTable) AddRow(dataRow *DataRow) {
	t.Rows = append(t.Rows, dataRow)
}

func (t *DataTable) AddError(index int64) {
	t.Errors = append(t.Errors, index)
}

func NewDataCell(key string, value string, rowNum int64, colNum int64, errs []string) *DataCell {
	cell := &DataCell{
		Key:    key,
		Value:  value,
		Col:    colNum,
		Row:    rowNum,
		Errors: errs,
	}
	return cell
}

func (r *DataRow) GetString(temp *Template, fieldKey string, setFunc func(v string), errFunc func(err ...error)) {
	f := temp.GetField(fieldKey)
	if f == nil && setFunc != nil {
		setFunc("")
		return
	}
	value, err := ValueToString(r.Values[fieldKey], "", f.GetReplace()...)
	if err != nil {
		errFunc(err)
	} else {
		setFunc(value)
	}
}

func (r *DataRow) GetFloat(temp *Template, fieldKey string, setFunc func(v *float64), errFunc func(err ...error)) {
	f := temp.GetField(fieldKey)
	if f == nil && setFunc != nil {
		setFunc(nil)
		return
	}
	value, err := ValueToFloat(r.Values[fieldKey], f.DefFloat)
	if err != nil {
		errFunc(err)
	} else {
		setFunc(value)
	}
}

func (r *DataRow) GetDate(temp *Template, fieldKey string, setFunc func(v *time.Time), errFunc func(err ...error)) {
	f := temp.GetField(fieldKey)
	if f == nil && setFunc != nil {
		setFunc(nil)
		return
	}
	value, err := ValueToDate(r.Values[fieldKey], f.GetDefTime(), f.GetReplace()...)
	if err != nil {
		errFunc(err)
	} else {
		setFunc(value)
	}
}

func (r *DataRow) GetBool(temp *Template, fieldKey string, defValue bool, setFunc func(v bool), errFunc func(err ...error)) {
	value, err := ValueToBool(r.Values[fieldKey], defValue)
	if err != nil {
		errFunc(err)
	} else {
		setFunc(value)
	}
}

func (r *DataRow) AddCell(fieldKey string, cell *DataCell) {
	if cell == nil {
		return
	}
	if cells, ok := r.Cells[fieldKey]; ok {
		r.Cells[fieldKey] = append(cells, cell)
	} else {
		r.Cells[fieldKey] = []*DataCell{cell}
	}
}

func (r *DataRow) AddValue(fieldKey string, value any) {
	r.Values[fieldKey] = value
}

func (r *DataRow) AddError(fieldKey string, err error) {
	if err == nil {
		return
	}
	if r.Errors == nil {
		r.Errors = map[string][]string{}
	}
	if errs, ok := r.Errors[fieldKey]; ok {
		r.Errors[fieldKey] = append(errs, err.Error())
	} else {
		r.Errors[fieldKey] = []string{err.Error()}
	}
}

func (r *DataRow) HasError() bool {
	return len(r.Errors) > 0
}

func ValueToFloat(value any, defValue *float64) (*float64, error) {
	var v any = value
	if v == nil {
		return defValue, nil
	}
	if i, ok := v.(goja.Value); ok {
		v = i.String()
	} else if f, ok := v.(float64); ok {
		return &f, nil
	} else if f, ok := v.(*float64); ok {
		return f, nil
	}
	if str, ok := v.(string); ok {
		if len(str) == 0 {
			return defValue, nil
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		return &f, err
	}

	return nil, errors.New("unsupported data type" + reflect.TypeOf(v).String())
}

func ValueToDate(v any, defValue *time.Time, replace ...Replace) (*time.Time, error) {
	var s = ""
	if i, ok := v.(goja.Value); ok {
		s, _ = ValueToString(i, "", replace...)
	}
	if value, ok := v.(time.Time); ok {
		return &value, nil
	}
	if value, ok := v.(*time.Time); ok {
		return value, nil
	}
	if str, ok := v.(string); ok {
		s = str
	}
	if len(s) == 0 {
		return defValue, nil
	}
	t, err := time.ParseInLocation(BaseFormat, s, setting.GetTimeZone())
	return &t, err
}

func ValueToInt(v any, defValue *int64) (*int64, error) {
	var s = ""
	if i, ok := v.(goja.Value); ok {
		s = i.String()
	}
	if value, ok := v.(int64); ok {
		return &value, nil
	}
	if value, ok := v.(*int64); ok {
		return value, nil
	}
	if str, ok := v.(string); ok {
		s = str
	}
	if len(s) == 0 {
		return defValue, nil
	}
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func ValueToBool(v any, defaultValue bool) (bool, error) {
	if i, ok := v.(goja.Value); ok {
		return i.ToBoolean(), nil
	}
	if value, ok := v.(bool); ok {
		return value, nil
	}
	if str, ok := v.(string); ok {
		if len(str) == 0 {
			return defaultValue, nil
		}
		f, err := strconv.ParseBool(str)
		return f, err
	}
	return defaultValue, nil
}

func ValueToString(v any, defValue string, replaces ...Replace) (res string, err error) {
	if v == nil {
		return defValue, nil
	}
	res = ""
	if i, ok := v.(goja.Value); ok {
		res = i.String()
	} else if value, ok := v.(string); ok {
		res = value
	} else if value, ok := v.(float64); ok {
		res, err = strconv.FormatFloat(value, 'f', 2, 64), nil
	} else if value, ok := v.(time.Time); ok {
		res = value.String()
	}

	for _, r := range replaces {
		res = strings.ReplaceAll(res, r.Old, r.New)
	}
	return res, nil
}
