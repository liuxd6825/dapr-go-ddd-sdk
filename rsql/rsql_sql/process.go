package rsql_sql

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"strings"
	"time"
)

type Process struct {
	sb       strings.Builder
	tenantId string
}

func NewProcess(tenantId string) rsql.Process {
	return &Process{tenantId: tenantId}
}

func (p *Process) TenantId() string {
	return p.tenantId
}

func (p *Process) GetSQL() string {
	return p.sb.String()
}

func (p *Process) GetFilter() any {
	return p.GetSQL()
}

func (p *Process) add(format string, a ...interface{}) {
	p.sb.WriteString(fmt.Sprintf(format, a...))
}

func (p *Process) OnFnProcess(expr rsql.Expression, fn *rsql.FuncValue) rsql.Value {
	switch fn.Name {
	case "sub":
		p := NewProcess(p.tenantId)
		err := rsql.ParseProcess(fn.RSQL(), p)
		if err != nil {
			panic(err)
		}
		sql := p.GetSQL()
		field := rsql.AsFieldName(fn.Args["field"])
		table := fn.Args["table"]
		sql = fmt.Sprintf("(select %s from %s where %s)", field, table, sql)
		return &rsql.StringValue{Value: sql}
	default:
		return fn
	}
	return fn
}

func (p *Process) OnEquals(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s=%v", name, val)
}

func (p *Process) OnNotEquals(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s!=(%v)", name, val)
}

func (p *Process) OnLike(name string, value any, rValue rsql.Value) {
	val := p.getLikeValue(rValue)
	p.add("%s like %v", name, val)
}

func (p *Process) OnNotLike(name string, value any, rValue rsql.Value) {
	val := p.getLikeValue(rValue)
	p.add("%s not like %v", name, val)
}

func (p *Process) OnContains(name string, value any, rValue rsql.Value) {
	if s, ok := rValue.(*rsql.StringValue); ok {
		p.add("%s like '%%%v%%'", name, s.Value)
	} else {
		panic("invalid rsql type in contains ")
	}
}

func (p *Process) OnNotContains(name string, value any, rValue rsql.Value) {
	if s, ok := rValue.(*rsql.StringValue); ok {
		p.add("%s not like '%%%v%%'", name, s.Value)
	} else {
		panic("invalid rsql type in OnNotContains ")
	}
}

func (p *Process) OnGreaterThan(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s>%v", name, val)
}

func (p *Process) OnGreaterThanOrEquals(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s>=%v", name, val)
}

func (p *Process) OnLessThan(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s<%v", name, val)
}

func (p *Process) OnLessThanOrEquals(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s<=%v", name, val)
}

func (p *Process) OnIn(name string, value any, rValue rsql.Value) {
	// 是标准数组查询
	if vList, ok := rValue.(*rsql.ListValue); ok {
		val := p.getInValue(vList)
		p.add("%s in (%s)", name, val)
	} else { // 是sub子查询
		val := p.getSubSql(rValue)
		p.add("%s in %s", name, val)
	}
}

func (p *Process) OnNotIn(name string, value any, rValue rsql.Value) {
	// 是标准数组查询
	if vList, ok := rValue.(*rsql.ListValue); ok {
		val := p.getInValue(vList)
		p.add("%s not in (%s)", name, val)
	} else { // 是sub子查询
		val := p.getSubSql(rValue)
		p.add("%s not in %s", name, val)
	}
}

func (p *Process) OnAndItem() {
	p.add(" and ")
}

func (p *Process) OnAndStart() {
	p.add("(")
}

func (p *Process) OnAndEnd() {
	p.add(")")
}

func (p *Process) OnOrItem() {
	p.add(" or ")
}

func (p *Process) OnOrStart() {
	p.add("(")
}

func (p *Process) OnOrEnd() {
	p.add(")")
}

func (p *Process) OnIsNull(name string, value any, rValue rsql.Value) {
	p.add("%s is null", name)
}

func (p *Process) OnNotIsNull(name string, value any, rValue rsql.Value) {
	p.add("%s is not null", name)
}

func (p *Process) OnStart(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s like '%s%%'", name, val)
}

func (p *Process) OnEnd(name string, value any, rValue rsql.Value) {
	val := p.getValue(rValue)
	p.add("%s like '%%%s'", name, val)
}

func (p *Process) getSubSql(rValue rsql.Value) string {
	val := p.getValue(rValue)
	if sub, ok := val.(string); ok {
		sub = strings.Trim(sub, "'")
		return sub
	}
	panic("invalid rsql type in sub sql")
}

func (p *Process) getValue(value rsql.Value) any {
	var v any
	var err error
	switch value.(type) {
	case *rsql.StringValue:
		sv, _ := value.(*rsql.StringValue)
		//v = sv.Value
		v = "'" + sv.Value + "'"
	case *rsql.IntegerValue:
		sv, _ := value.(*rsql.IntegerValue)
		v = sv.Value
	case *rsql.DateValue:
		sv, _ := value.(*rsql.DateValue)
		if date, er := time.Parse(rsql.DateLayout, sv.Value); er != nil {
			err = er
		} else {
			v = fmt.Sprintf("'%v-%v-%v'", date.Year(), date.Month(), date.Day())
		}

	case *rsql.DoubleValue:
		sv, _ := value.(*rsql.DoubleValue)
		v = sv.Value
	case *rsql.DateTimeValue:
		sv, _ := value.(*rsql.DateTimeValue)
		if date, er := time.Parse(rsql.DateTimeLayout, sv.Value); er != nil {
			err = er
		} else {
			v = fmt.Sprintf("'%v-%v-%v %v:%v:%v'", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
		}
	case *rsql.BooleanValue:
		sv, _ := value.(*rsql.BooleanValue)
		v = sv.Value
	case *rsql.ListValue:
		sv, _ := value.(*rsql.ListValue)
		v = p.getValueList(sv)
	case *rsql.FuncValue:
		sv, _ := value.(*rsql.FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}

func (p *Process) getLikeValue(value rsql.Value) string {
	var s = ""
	if strVal, ok := value.(*rsql.StringValue); ok {
		s = strVal.Value
	} else {
		s = fmt.Sprintf("%v", value)
	}
	s = strings.Replace(s, "'", "''", -1)
	s = strings.Replace(s, "*", "%", -1)
	return "'" + s + "'"
}

func (p *Process) getValueList(listValue *rsql.ListValue) []interface{} {
	list := make([]interface{}, 0)
	if listValue == nil {
		return list
	}
	for _, v := range listValue.Value {
		list = append(list, p.getValue(v))
	}
	return list
}

func (p *Process) getInValue(listValue *rsql.ListValue) string {
	sb := &strings.Builder{}
	count := len(listValue.Value)
	for i, v := range listValue.Value {
		if s, ok := v.(*rsql.StringValue); ok {
			sb.WriteString(fmt.Sprintf("'%v'", s))
		} else {
			sb.WriteString(fmt.Sprintf("%v", v))
		}
		if i < count-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}
