package mapper

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"time"
)

type JsonTimeConverter struct {
	jsonTime  types.JSONTime
	timeValue time.Time
}

func NewJsonTimeConverter() *JsonTimeConverter {
	jsonTime := types.JSONTime{}
	timeValue := time.Time{}
	return &JsonTimeConverter{
		jsonTime:  jsonTime,
		timeValue: timeValue,
	}
}

func (c *JsonTimeConverter) GetTypeConverters() []copier.TypeConverter {
	return []copier.TypeConverter{
		c.converter(c.jsonTime, c.timeValue),
		c.converter(c.jsonTime, &c.timeValue),
		c.converter(&c.jsonTime, c.timeValue),
		c.converter(&c.jsonTime, &c.timeValue),

		c.converter(c.timeValue, c.jsonTime),
		c.converter(c.timeValue, &c.jsonTime),
		c.converter(&c.timeValue, c.jsonTime),
		c.converter(&c.timeValue, &c.jsonTime),
	}
}

func (c *JsonTimeConverter) converter(srcType, dstType interface{}) copier.TypeConverter {
	converter := copier.TypeConverter{
		SrcType: srcType,
		DstType: dstType,
		Fn: func(src interface{}) (interface{}, error) {
			switch srcType.(type) {
			case *types.JSONDate:
				return c.pjsonToTime(src, srcType, dstType)
			case types.JSONDate:
				return c.jsonToTime(src, srcType, dstType)
			case *time.Time:
				return c.ptimeToJson(src, srcType, dstType)
			case time.Time:
				return c.timeToJson(src, srcType, dstType)
			}
			return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
		},
	}
	return converter
}

func (c *JsonTimeConverter) pjsonToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*types.JSONTime)
	if !ok {
		return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case time.Time:
		{
			return time.Time(*s), nil
		}
	case *time.Time:
		{
			res := time.Time(*s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.JsonTimeConverter() error: dst type not matching")
}

func (c *JsonTimeConverter) jsonToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(types.JSONTime)
	if !ok {
		return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case time.Time:
		{
			return time.Time(s), nil
		}
	case *time.Time:
		{
			res := time.Time(s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.JsonTimeConverter() error: dst type not matching")
}

func (c *JsonTimeConverter) ptimeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*time.Time)
	if !ok {
		return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case types.JSONTime:
		{
			return types.JSONTime(*s), nil
		}
	case *types.JSONTime:
		{
			res := types.JSONTime(*s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.JsonTimeConverter() error: dst type not matching")
}

func (c *JsonTimeConverter) timeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(time.Time)
	if !ok {
		return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case types.JSONTime:
		{
			return types.JSONTime(s), nil
		}
	case *types.JSONTime:
		{
			res := types.JSONTime(s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.JsonTimeConverter() error: dst type not matching")
}
