package maputils

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"strings"
	"time"
)

var stringType = reflect.TypeOf("time.Now()")
var dateType = reflect.TypeOf(time.Now())

const localTime = "2006-01-02 15:04:05"

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

func NewMap(structValue interface{}) (map[string]interface{}, error) {
	mapData := make(map[string]interface{})
	if err := Decode(structValue, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}

// hookFunc 数据类型转换hook
func hookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// 将string类型转换为time类型
		if !(f == stringType && t == dateType) {
			return data, nil
		}
		timeString := data.(string)
		theTime, err := time.Parse(time.RFC3339Nano, timeString)
		if err != nil {
			if theTime, err = time.Parse(localTime, timeString); err != nil {
				return data, err
			}
		}
		return theTime, err
	}
}

func decodeHook(fromType reflect.Type, toType reflect.Type, v interface{}) (interface{}, error) {
	if fromType.Kind() == reflect.String && toType.Name() == "Time" {
		return stringAsTime(v)
	}
	return v, nil
}

func stringAsTime(v interface{}) (time.Time, error) {
	sTime := v.(string)
	format := localTime
	if strings.Contains(sTime, "T") {
		format = time.RFC3339
	}
	res, err := time.Parse(format, sTime)
	if err != nil {
		res, err = time.Parse(time.RFC3339Nano, sTime)
	}
	return res, err
}
