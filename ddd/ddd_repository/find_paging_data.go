package ddd_repository

import "github.com/dapr/dapr-go-ddd-sdk/ddd"

//
// PagingData [T ddd.Entity]
// @Description:
// Parse 解析表达式
// or         : and ('OR' | 'or' and)*
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
//
type PagingData[T ddd.Entity] struct {
	Data       *[]T   `json:"data"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int64  `json:"totalPages"`
	PageNum    int64  `json:"pageNum"`
	PageSize   int64  `json:"pageSize"`
	Sort       string `json:"sort"`
	Filter     string `json:"filter"`
}

func NewPagingData[T ddd.Entity](data *[]T, totalRows int64, totalPages int64, filter string, sort string, pageSize int64, pageNum int64) *PagingData[T] {
	findData := &PagingData[T]{
		Data:       data,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Filter:     filter,
		Sort:       sort,
		PageSize:   pageSize,
		PageNum:    pageNum,
	}
	return findData
}

func (f *PagingData[T]) GetData() *[]T {
	return f.Data
}
