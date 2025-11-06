package times

import (
	"errors"
	"time"

	"github.com/jinzhu/copier"
)

type jsonDateConverter struct {
	jsonDate  Date
	timeValue time.Time
}

func NewJsonDateConverter() *jsonDateConverter {
	jsonData := Date{}
	timeValue := time.Time{}
	return &jsonDateConverter{
		jsonDate:  jsonData,
		timeValue: timeValue,
	}
}

func (c *jsonDateConverter) GetTypeConverters() []copier.TypeConverter {
	return []copier.TypeConverter{
		c.converter(c.jsonDate, c.timeValue),
		c.converter(c.jsonDate, &c.timeValue),
		c.converter(&c.jsonDate, c.timeValue),
		c.converter(&c.jsonDate, &c.timeValue),
	}
}

func (c *jsonDateConverter) converter(srcType, dstType interface{}) copier.TypeConverter {
	converter := copier.TypeConverter{
		SrcType: srcType,
		DstType: dstType,
		Fn: func(src interface{}) (interface{}, error) {
			switch srcType.(type) {
			case *Date:
				return c.pdateToTime(src, srcType, dstType)
			case Date:
				return c.dateToTime(src, srcType, dstType)
			case *time.Time:
				return c.ptimeToJson(src, srcType, dstType)
			case time.Time:
				return c.timeToJson(src, srcType, dstType)
			}
			return nil, errors.New("mapper.jsonDateConverter() error: src type not matching")
		},
	}
	return converter
}

func (c *jsonDateConverter) pdateToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*Date)
	if !ok {
		return nil, errors.New("mapper.jsonDateConverter() error: src type not matching")
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
	return nil, errors.New("mapper.jsonDateConverter() error: dst type not matching")
}

func (c *jsonDateConverter) dateToTime(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(Date)
	if !ok {
		return nil, errors.New("mapper.jsonDateConverter() error: src type not matching")
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
	return nil, errors.New("mapper.jsonDateConverter() error: dst type not matching")
}

func (c *jsonDateConverter) ptimeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(*time.Time)
	if !ok {
		return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case Date:
		{
			return Date(*s), nil
		}
	case *Date:
		{
			res := Date(*s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.jsonTimeConverter() error: dst type not matching")
}

func (c *jsonDateConverter) timeToJson(src, srcType, dstType interface{}) (interface{}, error) {
	s, ok := src.(time.Time)
	if !ok {
		return nil, errors.New("mapper.jsonTimeConverter() error: src type not matching")
	}
	switch dstType.(type) {
	case Date:
		{
			return Date(s), nil
		}
	case *Date:
		{
			res := Date(s)
			return &res, nil
		}
	}
	return nil, errors.New("mapper.jsonTimeConverter() error: dst type not matching")
}
