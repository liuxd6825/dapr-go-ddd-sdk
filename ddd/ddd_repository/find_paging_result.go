package ddd_repository

type PagingData struct {
	Data      interface{} `json:"data"`
	Count     int64       `json:"count"`
	Page      int64       `json:"page"`
	TotalPage int64       `json:"totalPage"`
	Size      int64       `json:"size"`
	Sort      string      `json:"sort"`
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
	Filter string `json:"filter"`
}

type FindPagingResult struct {
	FindResult
}

func NewFindPagingDataResult(data *PagingData, isFound bool, err error) *FindPagingResult {
	return &FindPagingResult{
		FindResult{
			data:    data,
			isFound: isFound,
			err:     err,
		},
	}
}

func (f *FindPagingResult) OnSuccess(success OnSuccessPaging) *FindPagingResult {
	if f.err == nil && success != nil && f.isFound {
		data := f.data.(*PagingData)
		f.err = success(data)
	}
	return f
}
