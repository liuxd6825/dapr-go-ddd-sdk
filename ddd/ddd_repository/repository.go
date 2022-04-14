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
	FindPaging(ctx context.Context, search *PagingQuery, opts ...*FindOptions) *FindPagingResult

	DoFind(fun func() (any, bool, error)) *FindResult
	DoSet(fun func() (interface{}, error)) *SetResult
	DoFilter(tenantId, filter string, fun func(filter map[string]interface{}) (*PagingData, bool, error)) *FindPagingResult
}

type Pageable struct {
	PageNum  int `json:"pageNum"`
	PageSize int `json:"pageSize"`
}

type PagingQuery struct {
	TenantId string `json:"tenantId"`
	Fields   string `json:"fields"`
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
	Filter   string `json:"filter"`
	Sort     string `json:"sort"` // name:desc,id:asc
	PageNum  int64  `json:"pageNum"`
	PageSize int64  `json:"pageSize"`
}

type Option struct {
}
