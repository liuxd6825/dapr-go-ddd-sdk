package rsql

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"time"
)

func GetValue(value Value) interface{} {
	var v interface{}
	switch value.(type) {
	case *StringValue:
		sv, _ := value.(*StringValue)
		v = sv.Value
	case *IntegerValue:
		sv, _ := value.(*IntegerValue)
		v = sv.Value
	case *DateValue:
		sv, _ := value.(*DateValue)
		v = sv.Value
	case *DoubleValue:
		sv, _ := value.(*DoubleValue)
		v = sv.Value
	case *DateTimeValue:
		sv, _ := value.(*DateTimeValue)
		v = sv.Value
	case *BooleanValue:
		sv, _ := value.(*BooleanValue)
		v = sv.Value
	case *ListValue:
		sv, _ := value.(*ListValue)
		v = GetValueList(sv)
	default:
		v = value
	}
	return v
}

func getValue(value Value) any {
	var v any
	var err error
	switch value.(type) {
	case *StringValue:
		sv, _ := value.(*StringValue)
		v = sv.Value
	case *IntegerValue:
		sv, _ := value.(*IntegerValue)
		v = sv.Value
	case *DateValue:
		sv, _ := value.(*DateValue)
		v, err = time.Parse(DateLayout, sv.Value)
	case *DoubleValue:
		sv, _ := value.(*DoubleValue)
		v = sv.Value
	case *DateTimeValue:
		sv, _ := value.(*DateTimeValue)
		v, err = time.Parse(DateTimeLayout, sv.Value)
	case *BooleanValue:
		sv, _ := value.(*BooleanValue)
		v = sv.Value
	case *ListValue:
		sv, _ := value.(*ListValue)
		v = GetValueList(sv)
	case *FuncValue:
		sv, _ := value.(*FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}

// AsFieldName
// @Description: 转换为mongodb规范的字段名称
// @param name
// @return string
func AsFieldName(name string) string {
	return stringutils.SnakeString(name)
}

func GetValueList(listValue *ListValue) []interface{} {
	list := make([]interface{}, 0)
	for _, v := range listValue.Value {
		list = append(list, GetValue(v))
	}
	return list
}
