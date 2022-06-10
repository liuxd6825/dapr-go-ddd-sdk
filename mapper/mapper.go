package mapper

import (
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"time"
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
	timeToString := copier.TypeConverter{
		SrcType: time.Time{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			v, err := timeToStr(src, types.TimeFormat)
			if err != nil {
				return nil, err
			}
			return types.TimeString(*v), err
		},
	}
	timeToTimeString := copier.TypeConverter{
		SrcType: time.Time{},
		DstType: types.TimeString(""),
		Fn: func(src interface{}) (interface{}, error) {
			v, err := timeToStr(src, types.TimeFormat)
			if err != nil {
				return nil, err
			}
			return types.TimeString(*v), err
		},
	}
	timeToDateString := copier.TypeConverter{
		SrcType: time.Time{},
		DstType: types.DateString(""),
		Fn: func(src interface{}) (interface{}, error) {
			v, err := timeToStr(src, types.DateFormat)
			if err != nil {
				return nil, err
			}
			return types.DateString(*v), err
		},
	}

	stringToTime := copier.TypeConverter{
		SrcType: copier.String,
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			return strToTime(src, types.TimeFormat)
		},
	}
	timeStringToTime := copier.TypeConverter{
		SrcType: types.TimeString(""),
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			return strToTime(src, types.TimeFormat)
		},
	}
	dateStringToTime := copier.TypeConverter{
		SrcType: types.DateString(""),
		DstType: time.Time{},
		Fn: func(src interface{}) (interface{}, error) {
			return strToTime(src, types.DateFormat)
		},
	}

	return []copier.TypeConverter{timeToString, timeToTimeString, timeToDateString, stringToTime, timeStringToTime, dateStringToTime}
}

func timeToStr(src interface{}, timeFormat string) (*string, error) {
	s, ok := src.(time.Time)
	if !ok {
		return nil, errors.New("mapper.timeToStr() error: src type not matching")
	}
	v := s.Format(timeFormat)
	return &v, nil
}

func strToTime(src interface{}, timeFormat string) (interface{}, error) {
	s, ok := src.(string)
	if !ok {
		return nil, errors.New("mapper.strToTime() error: src type not matching")
	}
	v, err := time.Parse(timeFormat, s)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("mapper.strToTime() error: string is \"%s\" ", s))
	}
	return v, err
}
