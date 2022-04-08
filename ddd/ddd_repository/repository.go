package ddd_repository

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
)

type Repository interface {
	Insert(ctx context.Context, entity ddd.Entity, opts ...*SetOptions) *SetResult
	Update(ctx context.Context, entity ddd.Entity, opts ...*SetOptions) *SetResult
	Delete(ctx context.Context, entity ddd.Entity, opts ...*SetOptions) *SetResult
	DeleteById(ctx context.Context, tenantId string, id string, opts ...*SetOptions) *SetResult
	FindById(ctx context.Context, tenantId string, id string, opts ...*FindOptions) *FindResult
	FindAll(ctx context.Context, tenantId string, opts ...*FindOptions) *FindResult
	FindPagingData(ctx context.Context, search *PagingQuery, opts ...*FindOptions) *FindPagingDataResult

	DoFind(fun func() (interface{}, bool, error)) *FindResult
	DoSet(fun func() (interface{}, error)) *SetResult
	DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*PagingData, bool, error)) *FindPagingDataResult
}

type Pageable struct {
	PageNumber int
	PageSize   int
}

type PagingQuery struct {
	TenantId string
	Fields   string
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
	Filter string
	Sort   string // name:desc,id:asc
	Page   int64
	Size   int64
}

type Option struct {
}
