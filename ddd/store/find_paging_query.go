package store

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
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
	SetMapToQuery(map[string]any) FindPagingQueryBuilder
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

func NewFindPagingQueryRequest() *FindPagingQueryRequest {
	return &FindPagingQueryRequest{}
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
	Cols []*GroupCol `json:"cols"`
}

type ValueCols struct {
	Cols []*ValueCol `json:"cols"`
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
		Cols: make([]*GroupCol, 0),
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
		groupCols.Cols = cols
	}
	return groupCols
}

func NewValueCols() *ValueCols {
	return &ValueCols{
		Cols: make([]*ValueCol, 0),
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
	return NewGroupKeysWidthString(s)
}

func NewGroupKeysWidthString(s string) []any {
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
	return NewGroupColsWidthString(s)
}

func NewGroupColsWidthString(s string) []*GroupCol {
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
	return NewValueColsWidth(s)
}

func NewValueColsWidth(s string) []*ValueCol {
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
	s.Cols = append(s.Cols, &GroupCol{Field: field, DataType: dataType})
	return s
}

func (s *GroupCols) GetCols() []*GroupCol {
	return s.Cols
}

func (s *ValueCols) Add(field string, aggFunc AggFunc) *ValueCols {
	s.Cols = append(s.Cols, &ValueCol{Field: field, AggFunc: aggFunc})
	return s
}

func (s *ValueCols) GetCols() []*ValueCol {
	return s.Cols
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

func (q *FindPagingQueryRequest) AsMap() map[string]any {
	data, err := maputils.NewMapJsonKey(q)
	if err != nil {
		panic(err)
	}
	return data
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

func (f *findPagingQueryBuilder) SetMapToQuery(m map[string]any) FindPagingQueryBuilder {

	if sort, err := maputils.GetString(m, "sort", ""); err == nil {
		f.query.SetSort(sort)
	} else {
		panic(err)
	}

	if tenantId, err := maputils.GetString(m, "tenantId", ""); err == nil {
		f.query.SetTenantId(tenantId)
	} else {
		panic(err)
	}

	if filter, err := maputils.GetString(m, "filter", ""); err == nil {
		f.query.SetFilter(filter)
	} else {
		panic(err)
	}

	if fields, err := maputils.GetString(m, "fields", ""); err == nil {
		f.query.SetFields(fields)
	} else {
		panic(err)
	}

	if pageNum, err := maputils.GetInt64(m, "pageNum", 0); err == nil {
		f.SetPageNum(pageNum)
	} else {
		panic(err)
	}

	if pageSize, err := maputils.GetInt64(m, "pageSize", 20); err == nil {
		f.SetPageSize(pageSize)
	} else {
		panic(err)
	}

	if v, err := maputils.GetString(m, "mustFilter", ""); err == nil {
		f.SetMustFilter(v)
	} else {
		panic(err)
	}

	if v, err := maputils.GetString(m, "groupCols", ""); err == nil {
		f.SetGroupCols(NewGroupColsWidthString(v))
	} else {
		panic(err)
	}

	if v, err := maputils.GetString(m, "groupKeys", ""); err == nil {
		f.SetGroupKeys(NewGroupKeysWidthString(v))
	} else {
		panic(err)
	}

	if v, err := maputils.GetString(m, "valueCols", ""); err == nil {
		f.SetValueCols(NewValueColsWidth(v))
	} else {
		panic(err)
	}

	return f
}

func (f *findPagingQueryBuilder) Build() FindPagingQuery {
	return f.query
}
