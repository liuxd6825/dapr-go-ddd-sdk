/*  Filter
// - name=="Kill Bill";year=gt=2003
// - name=="Kill Bill" and year>2003
// - genres=in=(sci-fi,action);(director=='Christopher Nolan',actor==*Bale);year=ge=2000
// - genres=in=(sci-fi,action) and (director=='Christopher Nolan' or actor==*Bale) and year>=2000
// - director.lastName==Nolan;year=ge=2000;year=lt=2010
// - director.lastName==Nolan and year>=2000 and year<2010
// - genres=in=(sci-fi,action);genres=out=(romance,animated,horror),director==Que*Tarantino
// - genres=in=(sci-fi,action) and genres=out=(romance,animated,horror) or director==Que*Tarantino
// or         : and ('OR' | 'or' and) *
// and        : constraint ('AND' | 'and' constraint)*
// constraint : group | comparison
// group      : '(' or ')'
// comparison : identifier comparator arguments
// identifier : [a-zA-Z0-9]+('.'[a-zA-Z0-9]+)*
// comparator : '==' | '!=' | '==~' | '!=~' | '>' | '>=' | '<' | '<=' | '=in=' | '=out='
// arguments  : '(' listValue ')' | value
// value      : int | double | string | date | datetime | boolean
// listValue  : value(','value)*
// int        : [0-9]+
// double     : [0-9]+'.'[0-9]*
// string     : '"'.*'"' | '\''.*'\''
// date       : [0-9]{4}'-'[0-9]{2}'-'\[0-9]{2}
// datetime   : date'T'[0-9]{2}':'[0-9]{2}':'[0-9]{2}('Z' | (('+'|'-')[0-9]{2}(':')?[0-9]{2}))?
// boolean    : 'true' | 'false'
*/

package ddd_repository

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"strings"
)

type FindPagingQuery interface {
	GetTenantId() string
	GetFields() string
	GetFilter() string
	GetSort() string
	GetPageNum() int64
	GetPageSize() int64
	GetIsTotalRows() bool
	GetGroupCols() []*GroupCol
	GetValueCols() []*ValueCol
	GetGroupKeys() []any

	SetTenantId(string)
	SetFields(string)
	SetFilter(string)
	SetSort(string)
	SetPageNum(int64)
	SetPageSize(int64)
	SetIsTotalRows(bool)
	SetGroupCols([]*GroupCol)
	SetValueCols([]*ValueCol)
	SetGroupKeys([]any)
}

type GroupCol struct {
	Field    string         `json:"field"`
	DataType types.DataType `json:"dataType"`
}

type AggFunc string

const (
	AggFuncSum   AggFunc = "sum"
	AggFuncCount AggFunc = "count"
	AggFuncAvg   AggFunc = "avg"
	AggFuncFirst AggFunc = "first"
	AggFuncLast  AggFunc = "last"
	AggFuncMax   AggFunc = "max"
	AggFuncMin   AggFunc = "min"
	AggFuncZero  AggFunc = "zero"
)

func (f AggFunc) Name() string {
	return string(f)
}

type ValueCol struct {
	AggFunc AggFunc `json:"aggFunc"`
	Field   string  `json:"field"`
}

type FindPagingQueryRequest struct {
	TenantId    string      `json:"tenantId"`
	Fields      string      `json:"fields"`
	Filter      string      `json:"filter"`
	Sort        string      `json:"sort"`
	PageNum     int64       `json:"pageNum"`
	PageSize    int64       `json:"pageSize"`
	IsTotalRows bool        `json:"isTotalRows"`
	GroupCols   []*GroupCol `json:"groupCols"`
	GroupKeys   []any       `json:"groupKeys"`
	ValueCols   []*ValueCol `json:"valueCols"`
}

type FindPagingQueryDTO struct {
	TenantId    string `json:"tenantId"`
	Fields      string `json:"fields"`
	Filter      string `json:"filter"`
	Sort        string `json:"sort"`
	PageNum     int64  `json:"pageNum"`
	PageSize    int64  `json:"pageSize"`
	IsTotalRows bool   `json:"isTotalRows"`
	GroupCols   string `json:"groupCols"`
	GroupKeys   string `json:"groupKeys"`
	ValueCols   string `json:"valueCols"`
}

func (d *FindPagingQueryDTO) NewFindPagingQueryRequest() *FindPagingQueryRequest {
	r := &FindPagingQueryRequest{}
	if d == nil {
		return r
	}
	r.PageNum = d.PageNum
	r.PageSize = d.PageSize
	r.Filter = d.Filter
	r.Fields = d.Fields
	r.TenantId = d.TenantId
	r.Sort = d.Sort
	r.IsTotalRows = d.IsTotalRows
	r.ValueCols = d.newValueCols(d.ValueCols)
	r.GroupCols = d.newGroupCols(d.GroupCols)
	r.GroupKeys = d.newGroupKeys(d.GroupKeys)
	return r
}

func (d *FindPagingQueryDTO) newGroupKeys(s string) []any {
	res := make([]any, 0)
	if len(s) == 0 {
		return res
	}
	list := strings.Split(s, ",")
	for _, key := range list {
		res = append(res, key)
	}
	return res
}

func (d *FindPagingQueryDTO) newGroupCols(s string) []*GroupCol {
	res := make([]*GroupCol, 0)
	maps := RSqlKeyValueToMap(s)
	for k, v := range maps {
		col := &GroupCol{
			Field:    k,
			DataType: types.DataType(v),
		}
		res = append(res, col)
	}
	return res
}

func (d *FindPagingQueryDTO) newValueCols(s string) []*ValueCol {
	res := make([]*ValueCol, 0)
	maps := RSqlKeyValueToMap(s)
	for k, v := range maps {
		col := &ValueCol{
			Field:   k,
			AggFunc: AggFunc(v),
		}
		res = append(res, col)
	}
	return res
}

type GroupCols struct {
	cols []*GroupCol
}

type ValueCols struct {
	cols []*ValueCol
}

func NewGroupColsByJson(jsonText string) ([]*GroupCol, error) {
	list := make([]*GroupCol, 0)
	if len(jsonText) > 0 {
		if err := json.Unmarshal([]byte(jsonText), &list); err != nil {
			return nil, err
		}

	}
	return list, nil
}

func NewValueColsByJson(jsonText string) ([]*ValueCol, error) {
	list := make([]*ValueCol, 0)
	if len(jsonText) > 0 {
		if err := json.Unmarshal([]byte(jsonText), &list); err != nil {
			return nil, err
		}

	}
	return list, nil
}

func NewGroupKeysByJson(jsonText string) ([]any, error) {
	list := make([]any, 0)
	if len(jsonText) > 0 {
		if err := json.Unmarshal([]byte(jsonText), &list); err != nil {
			return nil, err
		}

	}
	return list, nil
}

func NewFindPagingQuery() FindPagingQuery {
	query := &FindPagingQueryRequest{PageSize: 20}
	return query
}

func NewGroupCols() *GroupCols {
	return &GroupCols{
		cols: make([]*GroupCol, 0),
	}
}

func NewValueCols() *ValueCols {
	return &ValueCols{
		cols: make([]*ValueCol, 0),
	}
}

func (s *GroupCols) Add(field string, dataType types.DataType) *GroupCols {
	s.cols = append(s.cols, &GroupCol{Field: field, DataType: dataType})
	return s
}

func (s *GroupCols) Cols() []*GroupCol {
	return s.cols
}

func (s *ValueCols) Add(field string, aggFunc AggFunc) *ValueCols {
	s.cols = append(s.cols, &ValueCol{Field: field, AggFunc: aggFunc})
	return s
}

func (s *ValueCols) Cols() []*ValueCol {
	return s.cols
}

func (q *FindPagingQueryRequest) SetTenantId(value string) {
	q.TenantId = value
}

func (q *FindPagingQueryRequest) SetFields(value string) {
	q.Fields = value
}

func (q *FindPagingQueryRequest) SetFilter(value string) {
	q.Filter = value
}

func (q *FindPagingQueryRequest) SetSort(value string) {
	q.Sort = value
}

func (q *FindPagingQueryRequest) SetPageNum(value int64) {
	q.PageNum = value
}

func (q *FindPagingQueryRequest) SetPageSize(value int64) {
	q.PageSize = value
}

func (q *FindPagingQueryRequest) SetIsTotalRows(val bool) {
	q.IsTotalRows = val
}

func (q *FindPagingQueryRequest) SetGroupCols(value []*GroupCol) {
	q.GroupCols = value
}

func (q *FindPagingQueryRequest) SetValueCols(value []*ValueCol) {
	q.ValueCols = value
}

func (q *FindPagingQueryRequest) SetGroupKeys(val []any) {
	q.GroupKeys = val
}

func (q *FindPagingQueryRequest) GetTenantId() string {
	return q.TenantId
}

func (q *FindPagingQueryRequest) GetFields() string {
	return q.Fields
}

func (q *FindPagingQueryRequest) GetFilter() string {
	return q.Filter
}

func (q *FindPagingQueryRequest) GetSort() string {
	return q.Sort
}

func (q *FindPagingQueryRequest) GetPageNum() int64 {
	return q.PageNum
}

func (q *FindPagingQueryRequest) GetPageSize() int64 {
	return q.PageSize
}

func (q *FindPagingQueryRequest) GetIsTotalRows() bool {
	return q.IsTotalRows
}

func (q *FindPagingQueryRequest) GetValueCols() []*ValueCol {
	return q.ValueCols
}

func (q *FindPagingQueryRequest) GetGroupKeys() []any {
	return q.GroupKeys
}

func (q *FindPagingQueryRequest) GetGroupCols() []*GroupCol {
	return q.GroupCols
}

//
// Validate
// @Description: 命令数据验证
//
func (q *FindPagingQueryRequest) Validate() error {
	ve := errors.NewVerifyError()
	if len(q.TenantId) == 0 {
		ve.AppendField("TenantId", "不能为空")
	}
	return ve.GetError()
}
