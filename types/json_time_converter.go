package types

import (
	"errors"
	"github.com/jinzhu/copier"
	"time"
)

type jsonTimeConverter struct {
	jsonTime  JSONTime
	timeValue time.Time
}

func newJsonTimeConverter() *jsonTimeConverter {
	jsonTime := JSONTime{}
	timeValue := time.Time{}
	return &jsonTimeConverter{
		jsonTime:  jsonTime,
		timeValue: timeValue,
	}
}

func (c *jsonTimeConverter) getTypeConverters() []copier.TypeConverter {
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

func (c *jsonTimeConverter) converter(srcType, dstType interface{}) copier.TypeConverter {
	converter := copier.TypeConverter{
		SrcType: srcType,
		DstType: dstType,
		Fn: func(src interface{}) (interface{}, error) {
			switch srcType.(type) {
			case *JSONDate:
				return c.pjsonToTime(src, srcType, dstType)
			case JSONDate:
				return c.jsonToTime(src, srcType, dstType)
			case *time.Time:
				return c.ptimeToJson(src, srcType, dstType)
			case time.Time:
				return c.timeToJson(src, srcType, dstType)
			}
			return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
		},
	}
	return converter
}

func (c *jsonTimeConverter) pjsonToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*JSONTime)
	if !ok {
		return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
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
	return nil, errors.New("mapper.jsonTimeConverter() error: dst type not matching")
}

func (c *jsonTimeConverter) jsonToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(JSONTime)
	if !ok {
		return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
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
	return nil, errors.New("mapper.jsonTimeConverter() error: dst type not matching")
}

func (c *jsonTimeConverter) ptimeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*time.Time)
	if !ok {
		return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case JSONTime:
		{
			return JSONTime(*s), nil
		}
	case *JSONTime:
		{
			res := JSONTime(*s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.jsonTimeConverter() error: dst type not matching")
}

func (c *jsonTimeConverter) timeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(time.Time)
	if !ok {
		return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case JSONTime:
		{
			return JSONTime(s), nil
		}
	case *JSONTime:
		{
			res := JSONTime(s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.jsonTimeConverter() error: dst type not matching")
}
