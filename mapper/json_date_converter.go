package mapper

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"time"
)

type JsonDateConverter struct {
	jsonDate  types.JSONDate
	timeValue time.Time
}

func NewJsonDateConverter() *JsonDateConverter {
	jsonData := types.JSONDate{}
	timeValue := time.Time{}
	return &JsonDateConverter{
		jsonDate:  jsonData,
		timeValue: timeValue,
	}
}

func (c *JsonDateConverter) GetTypeConverters() []copier.TypeConverter {
	return []copier.TypeConverter{
		c.converter(c.jsonDate, c.timeValue),
		c.converter(c.jsonDate, &c.timeValue),
		c.converter(&c.jsonDate, c.timeValue),
		c.converter(&c.jsonDate, &c.timeValue),
	}
}

func (c *JsonDateConverter) converter(srcType, dstType interface{}) copier.TypeConverter {
	converter := copier.TypeConverter{
		SrcType: srcType,
		DstType: dstType,
		Fn: func(src interface{}) (interface{}, error) {
			switch srcType.(type) {
			case *types.JSONDate:
				return c.pdateToTime(src, srcType, dstType)
			case types.JSONDate:
				return c.dateToTime(src, srcType, dstType)
			case *time.Time:
				return c.ptimeToJson(src, srcType, dstType)
			case time.Time:
				return c.timeToJson(src, srcType, dstType)
			}
			return nil, errors.New("mapper.JsonDateConverter() error: src type not matching")
		},
	}
	return converter
}

func (c *JsonDateConverter) pdateToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*types.JSONDate)
	if !ok {
		return nil, errors.New("mapper.JsonDateConverter() error: src type not matching")
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
	return nil, errors.New("mapper.JsonDateConverter() error: dst type not matching")
}

func (c *JsonDateConverter) dateToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(types.JSONDate)
	if !ok {
		return nil, errors.New("mapper.JsonDateConverter() error: src type not matching")
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
	return nil, errors.New("mapper.JsonDateConverter() error: dst type not matching")
}

func (c *JsonDateConverter) ptimeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*time.Time)
	if !ok {
		return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case types.JSONDate:
		{
			return types.JSONDate(*s), nil
		}
	case *types.JSONDate:
		{
			res := types.JSONDate(*s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.JsonTimeConverter() error: dst type not matching")
}

func (c *JsonDateConverter) timeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(time.Time)
	if !ok {
		return nil, errors.New("mapper.JsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case types.JSONDate:
		{
			return types.JSONDate(s), nil
		}
	case *types.JSONDate:
		{
			res := types.JSONDate(s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.JsonTimeConverter() error: dst type not matching")
}
