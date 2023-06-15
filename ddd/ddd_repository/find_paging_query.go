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

import "github.com/liuxd6825/dapr-go-ddd-sdk/types"

type FindPagingQuery interface {
	GetTenantId() string
	GetFields() string
	GetFilter() string
	GetSort() string
	GetPageNum() int64
	GetPageSize() int64
	GetIsTotalRows() bool
	GetRowGroupCols() []*GroupCol
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

type FindPagingQueryObj struct {
	TenantId   string      `json:"tenantId"`
	Fields     string      `json:"fields"`
	Filter     string      `json:"filter"`
	Sort       string      `json:"sort"`
	PageNum    int64       `json:"pageNum"`
	PageSize   int64       `json:"pageSize"`
	IsTotalRow bool        `json:"isTotalRow"`
	GroupCols  []*GroupCol `json:"groupCols"`
	GroupKeys  []any       `json:"groupKeys"`
	ValueCols  []*ValueCol `json:"valueCols"`
}

type GroupCols struct {
	cols []*GroupCol
}

type ValueCols struct {
	cols []*ValueCol
}

func NewFindPagingQuery() FindPagingQuery {
	query := &FindPagingQueryObj{PageSize: 20}
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

func (q *FindPagingQueryObj) SetTenantId(value string) {
	q.TenantId = value
}

func (q *FindPagingQueryObj) SetFields(value string) {
	q.Fields = value
}

func (q *FindPagingQueryObj) SetFilter(value string) {
	q.Filter = value
}

func (q *FindPagingQueryObj) SetSort(value string) {
	q.Sort = value
}

func (q *FindPagingQueryObj) SetPageNum(value int64) {
	q.PageNum = value
}

func (q *FindPagingQueryObj) SetPageSize(value int64) {
	q.PageSize = value
}

func (q *FindPagingQueryObj) SetIsTotalRows(val bool) {
	q.IsTotalRow = val
}

func (q *FindPagingQueryObj) SetGroupCols(value []*GroupCol) {
	q.GroupCols = value
}

func (q *FindPagingQueryObj) SetValueCols(value []*ValueCol) {
	q.ValueCols = value
}

func (q *FindPagingQueryObj) SetGroupKeys(val []any) {
	q.GroupKeys = val
}

func (q *FindPagingQueryObj) GetTenantId() string {
	return q.TenantId
}

func (q *FindPagingQueryObj) GetFields() string {
	return q.Fields
}

func (q *FindPagingQueryObj) GetFilter() string {
	return q.Filter
}

func (q *FindPagingQueryObj) GetSort() string {
	return q.Sort
}

func (q *FindPagingQueryObj) GetPageNum() int64 {
	return q.PageNum
}

func (q *FindPagingQueryObj) GetPageSize() int64 {
	return q.PageSize
}

func (q *FindPagingQueryObj) GetIsTotalRows() bool {
	return q.IsTotalRow
}

func (q *FindPagingQueryObj) GetRowGroupCols() []*GroupCol {
	return q.GroupCols
}

func (q *FindPagingQueryObj) GetValueCols() []*ValueCol {
	return q.ValueCols
}

func (q *FindPagingQueryObj) GetGroupKeys() []any {
	return q.GroupKeys
}
