package mapper

import (
	"github.com/jinzhu/copier"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/mitchellh/mapstructure"
)

var option *copier.Option

func init() {
	option = getOption()
}

//
// Mapper
// @Description: 进行struct属性复制，支持深度复制
// @param fromObj 来源
// @param toObj 目标
// @return error
//
func Mapper(fromObj, toObj interface{}) error {
	return copier.CopyWithOption(toObj, fromObj, *option)
}

//
// MaskMapper
// @Description: 根据指定进行属性复制，不支持深度复制
// @param fromObj 来源
// @param toObj 目标
// @param mask 要复制属性列表
// @return error
//
func MaskMapper(fromObj, toObj interface{}, mask []string) error {
	var fromMap map[string]interface{}
	var err error
	switch fromObj.(type) {
	case *map[string]interface{}:
		value := fromObj.(*map[string]interface{})
		fromMap = *value
		break
	case map[string]interface{}:
		fromMap = fromObj.(map[string]interface{})
		break
	default:
		fromMap = make(map[string]interface{})
		if err = mapstructure.Decode(fromObj, &fromMap); err != nil {
			return err
		}
	}
	if len(mask) > 0 {
		maskMap := make(map[string]string)
		for _, key := range mask {
			maskMap[key] = stringutils.FirstUpper(key)
		}
		for key, _ := range fromMap {
			_, ok := maskMap[key]
			if !ok {
				delete(fromMap, key)
			}
		}
	}

	err = mapstructure.Decode(&fromMap, toObj)
	return err
}

func getOption() *copier.Option {
	return &copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters:  getTypeConverters(),
	}
}
func getTypeConverters() []copier.TypeConverter {
	var typeConverters []copier.TypeConverter
	typeConverters = append(typeConverters, NewJsonDateConverter().GetTypeConverters()...)
	typeConverters = append(typeConverters, NewJsonTimeConverter().GetTypeConverters()...)
	return typeConverters
}
