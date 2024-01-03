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
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"strings"
)

type FindPagingQuery interface {
	GetTenantId() string
	SetTenantId(string)

	GetFields() string
	SetFields(string)

	GetFilter() string
	SetFilter(string)

	GetMustFilter() string
	SetMustFilter(string)

	GetSort() string
	SetSort(string)

	GetPageNum() int64
	SetPageNum(int64)

	GetPageSize() int64
	SetPageSize(int64)

	GetIsTotalRows() bool
	SetIsTotalRows(bool)

	GetGroupCols() []*GroupCol
	SetGroupCols([]*GroupCol)

	GetValueCols() []*ValueCol
	SetValueCols([]*ValueCol)

	GetGroupKeys() []any
	SetGroupKeys([]any)
}

type FindPagingQueryBuilder interface {
	SetTenantId(string) FindPagingQueryBuilder
	SetFields(string) FindPagingQueryBuilder
	SetFilter(format string, value ...any) FindPagingQueryBuilder
	SetMustFilter(format string, value ...any) FindPagingQueryBuilder
	SetSort(format string, value ...any) FindPagingQueryBuilder
	SetPageNum(int64) FindPagingQueryBuilder
	SetPageSize(int64) FindPagingQueryBuilder
	SetIsTotalRows(bool) FindPagingQueryBuilder
	SetGroupCols([]*GroupCol) FindPagingQueryBuilder
	SetValueCols([]*ValueCol) FindPagingQueryBuilder
	SetGroupKeys([]any) FindPagingQueryBuilder
	Build() FindPagingQuery
}

type findPagingQueryBuilder struct {
	query FindPagingQuery
}

type GroupCol struct {
	Field    string         `json:"field"`
	DataType types.DataType `json:"dataType"`
}

type FindPagingQueryRequest struct {
	TenantId    string      `json:"tenantId"`
	Fields      string      `json:"fields"` // 以逗号分隔多个字段
	Filter      string      `json:"filter"`
	MustFilter  string      `json:"-"`
	Sort        string      `json:"sort"`
	PageNum     int64       `json:"pageNum"`
	PageSize    int64       `json:"pageSize"`
	IsTotalRows bool        `json:"isTotalRows"`
	GroupCols   []*GroupCol `json:"groupCols"`
	GroupKeys   []any       `json:"groupKeys"`
	ValueCols   []*ValueCol `json:"valueCols"`
}

type FindPagingQueryMustWhere interface {
	GetMustWhere() (string, error)
}

type FindPagingQueryDTO struct {
	TenantId    string `json:"tenantId"`
	Fields      string `json:"fields"`
	Filter      string `json:"filter"`
	MustFilter  string `json:"mustFilter"`
	Sort        string `json:"sort"`
	PageNum     int64  `json:"pageNum"`
	PageSize    int64  `json:"pageSize"`
	IsTotalRows bool   `json:"isTotalRows"`
	GroupCols   string `json:"groupCols"`
	GroupKeys   string `json:"groupKeys"`
	ValueCols   string `json:"valueCols"`
}

type FindByIdRequest struct {
	TenantId string `json:"tenantId"`
	Id       string `json:"id"`
}

type FindByAllRequest struct {
	TenantId string `json:"tenantId"`
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

type GroupCols struct {
	cols []*GroupCol
}

type ValueCols struct {
	cols []*ValueCol
}

func NewFindPagingQueryBuilder() FindPagingQueryBuilder {
	return &findPagingQueryBuilder{
		query: NewFindPagingQuery(),
	}
}

func NewFindPagingQuery() FindPagingQuery {
	query := &FindPagingQueryRequest{PageSize: 20}
	return query
}

func NewFindPagingQueryDTO() *FindPagingQueryDTO {
	return &FindPagingQueryDTO{}
}

func NewGroupCols(s string) *GroupCols {
	groupCols := &GroupCols{
		cols: make([]*GroupCol, 0),
	}
	if len(s) > 0 {
		cols := make([]*GroupCol, 0)
		maps := RSqlKeyValueToList(s)
		for _, v := range maps {
			col := &GroupCol{
				Field:    v.Key,
				DataType: types.DataType(v.Value),
			}
			cols = append(cols, col)
		}
		groupCols.cols = cols
	}
	return groupCols
}

func NewValueCols() *ValueCols {
	return &ValueCols{
		cols: make([]*ValueCol, 0),
	}
}

func (d *FindPagingQueryDTO) NewQuery() FindPagingQuery {
	return d.NewFindPagingQueryRequest()
}

func (d *FindPagingQueryDTO) NewFindPagingQueryRequest() *FindPagingQueryRequest {
	r := &FindPagingQueryRequest{}
	if d == nil {
		return r
	}
	r.PageNum = d.PageNum
	r.PageSize = d.PageSize
	r.Filter = d.Filter
	r.MustFilter = d.MustFilter
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
	maps := RSqlKeyValueToList(s)
	for _, v := range maps {
		col := &GroupCol{
			Field:    v.Key,
			DataType: types.DataType(v.Value),
		}
		res = append(res, col)
	}
	return res
}

func (d *FindPagingQueryDTO) newValueCols(s string) []*ValueCol {
	res := make([]*ValueCol, 0)
	maps := RSqlKeyValueToList(s)
	for _, v := range maps {
		col := &ValueCol{
			Field:   v.Key,
			AggFunc: AggFunc(v.Value),
		}
		res = append(res, col)
	}
	return res
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

func (q *FindPagingQueryRequest) GetTenantId() string {
	return q.TenantId
}

func (q *FindPagingQueryRequest) SetTenantId(value string) {
	q.TenantId = value
}

func (q *FindPagingQueryRequest) GetFields() string {
	return q.Fields
}

func (q *FindPagingQueryRequest) SetFields(value string) {
	q.Fields = value
}

func (q *FindPagingQueryRequest) GetFilter() string {
	return q.Filter
}

func (q *FindPagingQueryRequest) SetFilter(value string) {
	q.Filter = value
}

func (q *FindPagingQueryRequest) GetSort() string {
	return q.Sort
}

func (q *FindPagingQueryRequest) SetSort(value string) {
	q.Sort = value
}

func (q *FindPagingQueryRequest) GetPageNum() int64 {
	return q.PageNum
}
func (q *FindPagingQueryRequest) SetPageNum(value int64) {
	q.PageNum = value
}

func (q *FindPagingQueryRequest) GetPageSize() int64 {
	return q.PageSize
}

func (q *FindPagingQueryRequest) SetPageSize(value int64) {
	q.PageSize = value
}

func (q *FindPagingQueryRequest) GetMustFilter() string {
	return q.MustFilter
}

func (q *FindPagingQueryRequest) SetMustFilter(value string) {
	q.MustFilter = value
}

func (q *FindPagingQueryRequest) GetIsTotalRows() bool {
	return q.IsTotalRows
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

func (q *FindPagingQueryRequest) GetValueCols() []*ValueCol {
	return q.ValueCols
}

func (q *FindPagingQueryRequest) GetGroupKeys() []any {
	return q.GroupKeys
}

func (q *FindPagingQueryRequest) GetGroupCols() []*GroupCol {
	return q.GroupCols
}

// Validate
// @Description: 命令数据验证
func (q *FindPagingQueryRequest) Validate() error {
	ve := errors.NewVerifyError()
	if len(q.TenantId) == 0 {
		ve.AppendField("TenantId", "不能为空")
	}
	return ve.GetError()
}

func (r *FindByIdRequest) GetTenantId() string {
	return r.TenantId
}

func (r *FindByIdRequest) GetId() string {
	return r.Id
}

func (r *FindByAllRequest) GetTenantId() string {
	return r.TenantId
}

func (f *findPagingQueryBuilder) SetTenantId(s string) FindPagingQueryBuilder {
	f.query.SetTenantId(s)
	return f
}

func (f *findPagingQueryBuilder) SetFields(s string) FindPagingQueryBuilder {
	f.query.SetFields(s)
	return f
}

func (f *findPagingQueryBuilder) SetFilter(format string, value ...any) FindPagingQueryBuilder {
	f.query.SetFilter(fmt.Sprintf(format, value...))
	return f
}

func (f *findPagingQueryBuilder) SetMustFilter(format string, value ...any) FindPagingQueryBuilder {
	f.query.SetMustFilter(fmt.Sprintf(format, value...))
	return f
}

func (f *findPagingQueryBuilder) SetSort(format string, value ...any) FindPagingQueryBuilder {
	f.query.SetSort(fmt.Sprintf(format, value...))
	return f
}

func (f *findPagingQueryBuilder) SetPageNum(i int64) FindPagingQueryBuilder {
	f.query.SetPageNum(i)
	return f
}

func (f *findPagingQueryBuilder) SetPageSize(i int64) FindPagingQueryBuilder {
	f.query.SetPageSize(i)
	return f
}

func (f *findPagingQueryBuilder) SetIsTotalRows(b bool) FindPagingQueryBuilder {
	f.query.SetIsTotalRows(b)
	return f
}

func (f *findPagingQueryBuilder) SetGroupCols(cols []*GroupCol) FindPagingQueryBuilder {
	f.query.SetGroupCols(cols)
	return f
}

func (f *findPagingQueryBuilder) SetValueCols(cols []*ValueCol) FindPagingQueryBuilder {
	f.query.SetValueCols(cols)
	return f
}

func (f *findPagingQueryBuilder) SetGroupKeys(value []any) FindPagingQueryBuilder {
	f.query.SetGroupKeys(value)
	return f
}

func (f *findPagingQueryBuilder) Build() FindPagingQuery {
	return f.query
}
