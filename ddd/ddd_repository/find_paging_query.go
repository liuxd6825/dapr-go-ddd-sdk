package ddd_repository

type FindPagingQuery interface {
	TenantId() string
	Fields() string
	Filter() string
	Sort() string
	PageNum() int64
	PageSize() int64
}

func NewFindPagingQuery(tenantId, fields, filter, sort string, pageNum, pageSize int64) FindPagingQuery {
	return &findPagingQuery{
		tenantId: tenantId,
		filter:   filter,
		fields:   fields,
		sort:     sort,
		pageNum:  pageNum,
		pageSize: pageSize,
	}
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

func (q *findPagingQuery) TenantId() string {
	return q.tenantId
}

func (q *findPagingQuery) Fields() string {
	return q.fields
}

func (q *findPagingQuery) Filter() string {
	return q.filter
}

func (q *findPagingQuery) Sort() string {
	return q.sort
}

func (q *findPagingQuery) PageNum() int64 {
	return q.pageNum
}

func (q *findPagingQuery) PageSize() int64 {
	return q.pageSize
}
