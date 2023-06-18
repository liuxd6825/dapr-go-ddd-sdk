package maputils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"reflect"
	"time"
)

const (
	localTime        = "2006-01-02 15:04:05"
	timeTypeName     = "Time"
	jsonTimeTypeName = "JSONTime"
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

func Decode(input interface{}, out interface{}) error {
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
	if err := Decode(fromObj, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}

func NewFromStr(v []string) map[string]interface{} {
	mapData := make(map[string]interface{})
	for _, item := range v {
		mapData[item] = item
	}
	return mapData
}

func decodeHook(fromType reflect.Type, toType reflect.Type, v interface{}) (interface{}, error) {
	if fromType.Kind() == reflect.String && toType.Name() == timeTypeName {
		return timeutils.AnyToTime(v, time.Time{})
	} else if fromType.Kind() == reflect.String && toType.Name() == jsonTimeTypeName {
		return timeutils.AnyToTime(v, time.Time{})
	}
	return v, nil
}
