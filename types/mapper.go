package types

import (
	"github.com/jinzhu/copier"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"github.com/mitchellh/mapstructure"
)

type MaskType int

const (
	MaskTypeContain MaskType = iota // 包含字段
	MaskTypeExclude                 // 排除的字段
)

type MaskOptions struct {
	Mask []string
	Type MaskType
}

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

func MaskMapper(fromObj, toObj interface{}, mask []string) error {
	options := MaskOptions{
		Mask: mask,
		Type: MaskTypeContain,
	}
	return MaskMapperOptions(fromObj, toObj, &options)
}

func MaskMapperType(fromObj, toObj interface{}, mask []string, maskType MaskType) error {
	options := MaskOptions{
		Mask: mask,
		Type: maskType,
	}
	return MaskMapperOptions(fromObj, toObj, &options)
}

//
// MaskMapperOptions
// @Description: 根据指定进行属性复制，不支持深度复制
// @param fromObj 来源
// @param toObj 目标
// @param mask 要复制属性列表
// @return error
//
func MaskMapperOptions(fromObj, toObj interface{}, options *MaskOptions) error {
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
	if options != nil && len(options.Mask) > 0 {
		maskMap := make(map[string]string)
		for _, key := range options.Mask {
			name := stringutils.FirstUpper(key)
			maskMap[name] = name
		}
		for key, _ := range fromMap {
			_, ok := maskMap[key]
			maskType := options.Type
			switch maskType {
			case MaskTypeExclude:
				if ok {
					delete(fromMap, key)
				}
				break
			case MaskTypeContain:
				if !ok {
					delete(fromMap, key)
				}
				break
			}
		}
	}

	var metadata *mapstructure.Metadata
	config := &mapstructure.DecoderConfig{
		Result:   toObj,
		Metadata: metadata,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(fromMap)
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
	typeConverters = append(typeConverters, newJsonDateConverter().getTypeConverters()...)
	typeConverters = append(typeConverters, newJsonTimeConverter().getTypeConverters()...)
	return typeConverters
}
