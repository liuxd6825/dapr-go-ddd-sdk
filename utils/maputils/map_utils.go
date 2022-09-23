package maputils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"reflect"
	"time"
)

var stringType = reflect.TypeOf("time.Now()")
var dateType = reflect.TypeOf(time.Now())

const (
	localTime    = "2006-01-02 15:04:05"
	timeTypeName = "Time"
)

func GetString(m map[string]interface{}, key string, result *string, def string) (bool, error) {
	if v, ok := m[key]; ok {
		str := fmt.Sprintf("%v", v)
		result = &str
		return ok, nil
	}
	str := def
	result = &str
	return false, nil
}

func GetInt64(m map[string]interface{}, key string, result *int64, def int64) (bool, error) {
	if v, ok := m[key]; ok {
		str := v.(int64)
		result = &str
		return ok, nil
	}
	str := def
	result = &str
	return false, nil
}

func DecodeMap(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Squash:           true,
		Result:           out,
		DecodeHook:       decodeHook,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func NewMap(fromObj interface{}) (map[string]interface{}, error) {
	mapData := make(map[string]interface{})
	if err := DecodeMap(fromObj, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}

func decodeHook(fromType reflect.Type, toType reflect.Type, v interface{}) (interface{}, error) {
	println("formType:" + fromType.Elem().Name())
	println("toType:" + toType.Elem().Name())
	if fromType.Kind() == reflect.String && toType.Name() == timeTypeName {
		return timeutils.AnyToTime(v, time.Time{})
	}
	return v, nil
}
