package engine

import (
	"bytes"
	"encoding/json"
	"strings"

	sch "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/liuxd6825/jsonschema/v6"
)

func IfElse(condition bool, trueVal any, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}

func ToJsonString(sch *jsonschema.Schema) string {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(true) // 禁用HTML转义
	err := encoder.Encode(sch.GetView())
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(buf.String(), "${", "&#36;{")
}

func NullQuery(query *sch.Query) bool {
	if query == nil {
		return true
	}
	if len(*query) > 0 {
		return false
	}
	return true
}

func NullForm(form *sch.Form) bool {
	if form == nil {
		return true
	}
	if len(*form) > 0 {
		return false
	}
	return true
}

func FormValue(form *sch.Form, formName string, propName string) any {
	if form == nil {
		return nil
	}

	if formName == "" {
		if val, ok := (*form)[propName]; ok {
			return val
		}
		return nil
	}

	formVal, ok := (*form)[formName]
	if ok {
		if propName == "" {
			return formVal
		} else {
			return MapValue(formVal.(map[string]any), propName)
		}
	}
	return nil
}

func NullColumn(column *sch.Column) bool {
	if column == nil {
		return true
	}
	if len(*column) > 0 {
		return false
	}
	return true
}

func ColumnValue(column *sch.Column, gridName string, propName string) any {
	if column == nil {
		return nil
	}

	if gridName == "" {
		if val, ok := (*column)[propName]; ok {
			return val
		}
		return nil
	}

	val, ok := (*column)[gridName]
	if ok {
		if propName == "" {
			return val
		} else {
			return MapValue(val.(map[string]any), propName)
		}
	}
	return nil
}

func MapValue(param map[string]any, propName string) any {
	if param == nil {
		return nil
	}
	if val, ok := param[propName]; ok {
		return val
	}
	return nil
}

func OnlyField(param map[string]any) bool {
	return OnlyOneField(param, "field")
}

func OnlyOneField(param map[string]any, fieldName string) bool {
	if param != nil && len(param) == 1 {
		if _, ok := param[fieldName]; ok {
			return true
		}
	}
	return false
}

func PropValue(field string, props []*Property) *Property {
	for _, prop := range props {
		if prop.Name() == field {
			return prop
		}
	}
	return nil
}

var globalValues *types.CMap[any] = types.NewCMap[any]()

func GetGlobal(name string) any {
	val, ok := globalValues.Get(name)
	if ok {
		return val
	}
	return nil
}

func SetGlobal(name string, val any) {
	globalValues.Set(name, val)
}
