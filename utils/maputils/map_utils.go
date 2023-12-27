package maputils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"time"
)

type JsonTime interface {
	PTime() *time.Time
	Time() time.Time
}

const (
	localTime        = "2006-01-02 15:04:05"
	timeTypeName     = "Time"
	jsonTimeTypeName = "JSONTime"
)

func MapToSnakeKey(data map[string]any) map[string]any {
	m := make(map[string]any)
	for k, v := range data {
		key := stringutils.SnakeString(k)
		if mv, ok := v.(map[string]any); ok {
			m[key] = MapToSnakeKey(mv)
		} else if mv, ok := v.(bson.M); ok {
			m[key] = MapToSnakeKey(mv)
		} else {
			m[key] = v
		}
	}
	return m
}

func GetKeys(data map[string]any) []string {
	keys := make([]string, len(data))
	i := 0
	for key, _ := range data {
		keys[i] = key
		i++
	}
	return keys
}

func GetKeysToFirstLower(data map[string]any) []string {
	keys := make([]string, len(data))
	i := 0
	for key, _ := range data {
		keys[i] = stringutils.FirstLower(key)
		i++
	}
	return keys
}

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
func NewMap(fromObj any) (map[string]any, error) {
	mapData := make(map[string]any)
	if err := Decode(fromObj, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}
func NewMapSnakeKey(fromObj any) (map[string]any, error) {
	return NewMapWithOptions(fromObj, nil, true)
}

func NewMapJsonKey(fromObj any) (map[string]any, error) {
	return NewMapWith(fromObj, nil, false, true)
}

func NewMapWithOptions(fromObj any, mask []string, snakeKey bool) (map[string]any, error) {
	return NewMapWith(fromObj, mask, snakeKey, false)
}

func NewMapWith(fromObj any, mask []string, snakeKey bool, jsonKey bool) (map[string]any, error) {
	mapData := make(map[string]any)
	if err := Decode(fromObj, &mapData); err != nil {
		return nil, err
	}

	if len(mask) > 0 {
		maskMap := map[string]any{}
		for _, key := range mask {
			key = stringutils.CamelString(key)
			maskMap[key] = mapData[key]
		}
		mapData = maskMap
	}

	if snakeKey || jsonKey {
		data := make(map[string]any)
		for k, v := range mapData {
			key := k
			if snakeKey {
				key = stringutils.SnakeString(key)
			}
			if jsonKey {
				key = stringutils.FirstLower(key)
			}
			data[key] = v
		}
		return data, nil
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
