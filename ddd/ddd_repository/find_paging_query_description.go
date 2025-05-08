package ddd_repository

const FilterDescription = `
- name=="Bill";year=gt=2003
- name=="Bill" and (year>2003 or year==2000)
- genres=in=(sci-fi,action);(director=='Christopher Nolan',actor==*Bale);year=ge=2000
- genres=in=(sci-fi,action) and (director=='Christopher Nolan' or actor==*Bale) and year>=2000
- director.lastName==Nolan;year=ge=2000;year=lt=2010
- director.lastName==Nolan and year>=2000 and year<2010
- genres=in=(sci-fi,action);genres=out=(romance,animated,horror),director==Que*Tarantino
- genres=in=(sci-fi,action) and genres=out=(romance,animated,horror) or director==Que*Tarantino
 or         : and ('OR' | 'or' and) *
 and        : constraint ('AND' | 'and' constraint)*
 constraint : group | comparison
 group      : '(' or ')'
 comparison : identifier comparator arguments
 identifier : [a-zA-Z0-9]+('.'[a-zA-Z0-9]+)*
 comparator : '==' | '!=' | '==~' | '!=~' | '>' | '>=' | '<' | '<=' | '=in=' | '=out='
 arguments  : '(' listValue ')' | value
 value      : int | double | string | date | datetime | boolean
 listValue  : value(','value)*
 int        : [0-9]+
 double     : [0-9]+'.'[0-9]*
 string     : '"'.*'"' | '\''.*'\''
 date       : [0-9]{4}'-'[0-9]{2}'-'\[0-9]{2}
 datetime   : date'T'[0-9]{2}':'[0-9]{2}':'[0-9]{2}('Z' | (('+'|'-')[0-9]{2}(':')?[0-9]{2}))?
 boolean    : 'true' | 'false'
`

const SortDescription = `field1:desc, field2:asc, ....`

const FieldsDescription = `field1, field2, field3, ...`

const ValueColsDescription = `field1:$func, field:$func, ...
payout:sum,income:count,amount:avg
$func = sum | count | avg | first | last | max | min | zer 
`

const GroupColsDescription = `field1:$type, field2:$type, ...
$type = string | int | float | money | date | dateTime | bool | array | object | year | month | day 
`

const GroupKeysDescription = `field1, field2, field3, ...`
