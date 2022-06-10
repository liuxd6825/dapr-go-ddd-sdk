package mapper

import (
	"github.com/jinzhu/copier"
)

var option *copier.Option

func init() {
	option = getOption()
}

func Mapper(fromObj, toObj interface{}) error {
	return copier.CopyWithOption(toObj, fromObj, *option)
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
