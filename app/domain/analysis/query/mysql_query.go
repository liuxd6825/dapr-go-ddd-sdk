package query

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

type AutoCompleteQuery struct {
	Field  string `json:"field" query:"field"`
	Value  string `json:"value" query:"value"`
	Filter string `json:"filter" query:"filter"`
}

type MysqlQuery struct {
	Params     map[string]MysqlQueryParam `json:"params" title:"参数"`
	Collection string
	QuerySQL   string
	Pipeline   []map[string]any `json:"pipeline" validate:"required" title:"管道"`
	Options    Options
}

type MysqlQueryParam struct {
	Value      string `json:"value"`
	DataType   string `json:"dataType"`
	NullRemove bool   `json:"nullRemove"`
	Formatter  string `json:"formatter"`
}

type Options struct {
	Sort *SortOptions
}

type SortOptions struct {
	TimeKey  string
	TimeType model.SortTimeType `json:"timeType" title:"排序类型"`
}

func (q *MysqlQuery) Init() {
	for _, v := range q.Pipeline {
		if matchVal, ok := v["$match"]; ok {
			if match, ok := matchVal.(map[string]any); ok {
				v["$match"] = q.initMatch(match)
			}
		}
	}
	return
}

func (q *MysqlQuery) GetMatch() map[string]any {
	for _, v := range q.Pipeline {
		if matchVal, ok := v["$match"]; ok {
			if match, ok := matchVal.(map[string]any); ok {
				return match
			}
		}
	}
	return map[string]any{}
}

func (p *MysqlQueryParam) getValue() (res any, err error) {
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

func (p *MysqlQueryParam) isNull() bool {
	if p.Value == "" {
		return true
	}
	return false
}

func (q *MysqlQuery) initMatch(match map[string]any) map[string]any {
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

func (q *MysqlQuery) getFieldValue(values map[string]any) (res map[string]any, remove bool, err error) {
	res = make(map[string]any)
	for key, val := range values {
		switch v := val.(type) {
		case map[string]interface{}:
			valMap := val.(map[string]any)
			res[key] = q.initMatch(valMap)
		case []any:
			processed := q.processSliceValues(v)
			if len(processed) > 0 {
				res[key] = processed
			}
		case int64, int32:
			res[key] = val
		case string:
			isRemove := false
			val, isRemove, err = q.getValue(v)
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

func (q *MysqlQuery) processSliceValues(items []any) []any {
	res := make([]any, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok {
			processed, isRemove, err := q.getValue(s)
			if err != nil {
				panic(err)
			}
			if isRemove {
				continue
			}
			res = append(res, processed)
			continue
		}
		res = append(res, item)
	}
	return res
}

func (q *MysqlQuery) getValue(val string) (res any, remove bool, err error) {
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

func (q *MysqlQuery) BuildWhereClause(match map[string]any) string {
	if len(match) == 0 {
		return "1=1"
	}

	var conditions []string
	dateRange := q.buildDateRangeCondition()
	if dateRange != "" {
		conditions = append(conditions, dateRange)
	}

	for key, val := range match {
		if key == "date" {
			continue
		}
		switch v := val.(type) {
		case map[string]any:
			for op, value := range v {
				switch op {
				case "$gte":
					conditions = append(conditions, fmt.Sprintf("%s >= '%v'", key, value))
				case "$lte":
					conditions = append(conditions, fmt.Sprintf("%s <= '%v'", key, value))
				case "$gt":
					conditions = append(conditions, fmt.Sprintf("%s > '%v'", key, value))
				case "$lt":
					conditions = append(conditions, fmt.Sprintf("%s < '%v'", key, value))
				case "$eq":
					conditions = append(conditions, fmt.Sprintf("%s = '%v'", key, value))
				case "$ne":
					conditions = append(conditions, fmt.Sprintf("%s != '%v'", key, value))
				case "$in":
					conditions = append(conditions, buildInClause(key, value))
				}
			}
		case string:
			conditions = append(conditions, fmt.Sprintf("%s = '%v'", key, v))
		case int, int32, int64:
			conditions = append(conditions, fmt.Sprintf("%s = %v", key, v))
		}
	}

	result := "1=1"
	if len(conditions) > 0 {
		for i, c := range conditions {
			if i == 0 {
				result = c
			} else {
				result += " AND " + c
			}
		}
	}
	return result
}

func buildInClause(key string, value any) string {
	items, ok := value.([]any)
	if !ok || len(items) == 0 {
		return "1=0"
	}
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		quoted = append(quoted, fmt.Sprintf("'%v'", item))
	}
	return fmt.Sprintf("%s IN (%s)", key, strings.Join(quoted, ","))
}

func (q *MysqlQuery) buildDateRangeCondition() string {
	startYear, ok1 := q.Params["startYear"]
	endYear, ok2 := q.Params["endYear"]

	if !ok1 || !ok2 {
		return ""
	}

	startVal, _ := startYear.getValue()
	endVal, _ := endYear.getValue()

	if startVal == nil || endVal == nil {
		return ""
	}

	var startStr, endStr string
	if t, ok := startVal.(time.Time); ok {
		startStr = t.Format("2006-01-02")
	} else {
		startStr = fmt.Sprintf("%v", startVal)
	}
	if t, ok := endVal.(time.Time); ok {
		endStr = t.Format("2006-01-02")
	} else {
		endStr = fmt.Sprintf("%v", endVal)
	}

	return fmt.Sprintf("date >= '%s' AND date <= '%s 23:59:59'", startStr, endStr)
}

func (q *MysqlQuery) BuildGroupClause() string {
	if q.Options.Sort == nil {
		return "year, month"
	}
	switch q.Options.Sort.TimeType {
	case model.Year:
		return "year"
	case model.Month:
		return "year, month"
	case model.Day:
		return "year, month, day"
	default:
		return "year, month"
	}
}

func (q *MysqlQuery) BuildSelectClause() string {
	if q.Options.Sort == nil {
		return `DATE_FORMAT(CONCAT(year, '-', LPAD(month, 2, '0'), '-01'), '%Y-%m-%d') as time,
			SUM(amount) as amount,
			COUNT(*) as record,
			COUNT(DISTINCT opp_name) as oppCount`
	}
	switch q.Options.Sort.TimeType {
	case model.Year:
		return `DATE_FORMAT(CONCAT(year, '-01-01'), '%Y-%m-%d') as time,
			SUM(amount) as amount,
			COUNT(*) as record,
			COUNT(DISTINCT opp_name) as oppCount`
	case model.Month:
		return `DATE_FORMAT(CONCAT(year, '-', LPAD(month, 2, '0'), '-01'), '%Y-%m-%d') as time,
			SUM(amount) as amount,
			COUNT(*) as record,
			COUNT(DISTINCT opp_name) as oppCount`
	case model.Day:
		return `DATE_FORMAT(CONCAT(year, '-', LPAD(month, 2, '0'), '-', LPAD(day, 2, '0')), '%Y-%m-%d') as time,
			SUM(amount) as amount,
			COUNT(*) as record,
			COUNT(DISTINCT opp_name) as oppCount`
	default:
		return `DATE_FORMAT(CONCAT(year, '-', LPAD(month, 2, '0'), '-01'), '%Y-%m-%d') as time,
			SUM(amount) as amount,
			COUNT(*) as record,
			COUNT(DISTINCT opp_name) as oppCount`
	}
}

func (q *MysqlQuery) BuildOrderClause() string {
	return "time ASC"
}
