package ddd_repository

type FindPagingQuery interface {
	GetTenantId() string
	GetFields() string
	GetFilter() string
	GetSort() string
	GetPageNum() int64
	GetPageSize() int64
}

type FindPagingQueryOptions struct {
	tenantId *string
	fields   *string
	filter   *string
	sort     *string
	pageNum  *int64
	pageSize *int64
}

func NewFindPagingQueryOptions() *FindPagingQueryOptions {
	return &FindPagingQueryOptions{}
}

func NewFindPagingQueryOptionsAll(fields, filter, sort string, pageNum, pageSize int64) *FindPagingQueryOptions {
	return &FindPagingQueryOptions{
		filter:   &filter,
		fields:   &fields,
		sort:     &sort,
		pageNum:  &pageNum,
		pageSize: &pageSize,
	}
}

func NewFindPagingQueryOptionsDefault() *FindPagingQueryOptions {
	fields := ""
	filter := ""
	sort := ""
	pageNum := int64(0)
	pageSize := int64(20)
	return &FindPagingQueryOptions{
		fields:   &fields,
		filter:   &filter,
		sort:     &sort,
		pageNum:  &pageNum,
		pageSize: &pageSize,
	}
}

func (o *FindPagingQueryOptions) SetFields(fields string) *FindPagingQueryOptions {
	o.fields = &fields
	return o
}

func (o *FindPagingQueryOptions) SetFilter(filter string) *FindPagingQueryOptions {
	o.filter = &filter
	return o
}

func (o *FindPagingQueryOptions) SetSort(sort string) *FindPagingQueryOptions {
	o.sort = &sort
	return o
}

func (o *FindPagingQueryOptions) SetPageNum(pageNum int64) *FindPagingQueryOptions {
	o.pageNum = &pageNum
	return o
}

func (o *FindPagingQueryOptions) SetpPageSize(pageSize int64) *FindPagingQueryOptions {
	o.pageSize = &pageSize
	return o
}

func NewFindPagingQuery(tenantId string, options ...*FindPagingQueryOptions) FindPagingQuery {
	opt := NewFindPagingQueryOptionsDefault()
	if options != nil {
		for _, o := range options {
			if o.sort != nil {
				opt.sort = o.sort
			}
			if o.pageNum != nil {
				opt.pageNum = o.pageNum
			}
			if o.pageSize != nil {
				opt.pageSize = o.pageSize
			}
			if o.fields != nil {
				opt.fields = o.fields
			}
			if o.filter != nil {
				opt.filter = o.filter
			}
		}
	}

	query := &findPagingQuery{
		tenantId: tenantId,
		filter:   *opt.filter,
		fields:   *opt.fields,
		sort:     *opt.sort,
		pageNum:  *opt.pageNum,
		pageSize: *opt.pageSize,
	}
	return query
}

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
type findPagingQuery struct {
	tenantId string
	fields   string
	filter   string
	sort     string
	pageNum  int64
	pageSize int64
}

func (q *findPagingQuery) GetTenantId() string {
	return q.tenantId
}

func (q *findPagingQuery) GetFields() string {
	return q.fields
}

func (q *findPagingQuery) GetFilter() string {
	return q.filter
}

func (q *findPagingQuery) GetSort() string {
	return q.sort
}

func (q *findPagingQuery) GetPageNum() int64 {
	return q.pageNum
}

func (q *findPagingQuery) GetPageSize() int64 {
	return q.pageSize
}
