package query

import (
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

type AggregateQuery struct {
	Params     map[string]AggregateParam `json:"params" title:"参数"`
	Collection string                    `json:"collection" validate:"required" title:"集合名称"`
	Pipeline   []map[string]any          `json:"pipeline" validate:"required" title:"管道"`
	Options    AggregateOptions          `json:"options" title:"选项"`
}

type AggregateOptions struct {
	Sort *MongoQuerySort `json:"sort" title:"补齐排序"`
}

type AggregateParam struct {
	Value      string `json:"value"`
	DataType   string `json:"dataType"`
	NullRemove bool   `json:"nullRemove"`
	Formatter  string `json:"formatter"`
}

type MongoQuerySort struct {
	TimeKey  string             `json:"timeKey"  title:"时间字段"`
	TimeType model.SortTimeType `json:"timeType" title:"排序类型"`
}

func (q *AggregateQuery) Init() {
	for _, v := range q.Pipeline {
		if matchVal, ok := v["$match"]; ok {
			if match, ok := matchVal.(map[string]any); ok {
				v["$match"] = q.initMatch(match)
			}
		}
	}
	return
}

func (q *AggregateQuery) GetMatch() map[string]any {
	for _, v := range q.Pipeline {
		if matchVal, ok := v["$match"]; ok {
			if match, ok := matchVal.(map[string]any); ok {
				return match
			}
		}
	}
	return map[string]any{}
}

func (p *AggregateParam) getValue() (res any, err error) {
	val := p.Value
	if p.Formatter != "" {
		val = fmt.Sprintf(p.Formatter, p.Value)
	}
	res = val
	if p.DataType == "date" {
		if !p.isNull() {
			res, err = timeutils.AsTime(val)
		} else {
			if p.NullRemove {
				res = nil
			} else {
				err = errors.New("")
			}
		}
	}
	return res, err
}

func (p *AggregateParam) isNull() bool {
	if p.Value == "" {
		return true
	}
	return false
}

func (q *AggregateQuery) initMatch(match map[string]any) map[string]any {
	m := make(map[string]any)
	for key, val := range match {
		switch val.(type) {
		case map[string]interface{}:
			valMap := val.(map[string]any)
			valData, remove, err := q.getFieldValue(valMap)
			if err != nil {
				panic(err)
			}
			if !remove {
				m[key] = valData
			}
		case int64, int32:
			m[key] = val
		case string:
			val, remove, err := q.getValue(val.(string))
			if err != nil {
				panic(err)
			}
			if !remove {
				m[key] = val
			}
		}
	}
	return m
}

func (q *AggregateQuery) getFieldValue(values map[string]any) (res map[string]any, remove bool, err error) {
	res = make(map[string]any)
	for key, val := range values {
		switch val.(type) {
		case map[string]interface{}:
			valMap := val.(map[string]any)
			res[key] = q.initMatch(valMap)
		case int64, int32:
			res[key] = val
		case string:
			isRemove := false
			val, isRemove, err = q.getValue(val.(string))
			if err != nil {
				panic(err)
			}
			if !isRemove {
				res[key] = val
			}
		}
	}
	return res, len(res) == 0, err
}

func (q *AggregateQuery) getValue(val string) (res any, remove bool, err error) {
	res = val
	for name, param := range q.Params {
		if val == "@"+name {
			if param.isNull() && param.NullRemove {
				remove = true
				return nil, true, nil
			}
			res, err = param.getValue()
		}
	}
	return res, remove, err
}
