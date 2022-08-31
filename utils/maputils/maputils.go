package maputils

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

var stringType = reflect.TypeOf("time.Now()")
var dateType = reflect.TypeOf(time.Now())

const localTime = "2006-01-02 15:04:05"

func Decode(input interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Squash:           true,
		Result:           out,
		DecodeHook:       hookFunc(),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
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
