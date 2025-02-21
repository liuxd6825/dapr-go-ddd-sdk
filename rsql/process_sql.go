package rsql

import (
	"fmt"
	"strings"
	"time"
)

type sqlProcess struct {
	str      string
	tenantId string
}

func NewSqlProcess(tenantId string) Process {
	return &sqlProcess{tenantId: tenantId}
}

func (p *sqlProcess) TenantId() string {
	return p.tenantId
}

func (p *sqlProcess) OnFnProcess(expr Expression, fn *FuncValue) Value {
	switch fn.Name {
	case "sub":
		p := NewSqlProcess(p.tenantId)
		err := ParseProcess(fn.RSQL(), p)
		if err != nil {
			panic(err)
		}
		sql := p.GetSQL()
		field := AsFieldName(fn.Args["field"])
		table := fn.Args["table"]
		sql = fmt.Sprintf("(select %s from %s where %s)", field, table, sql)
		return &StringValue{Value: sql}
	default:
		return fn
	}
	return fn
}

func (p *sqlProcess) OnEquals(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, val)
}

func (p *sqlProcess) OnNotEquals(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s != (%v)", p.str, name, val)
}

func (p *sqlProcess) OnLike(name string, value any, rValue Value) {
	val := p.getLikeValue(rValue)
	p.str = fmt.Sprintf("%s %s like %v", p.str, name, val)
}

func (p *sqlProcess) OnNotLike(name string, value any, rValue Value) {
	val := p.getLikeValue(rValue)
	p.str = fmt.Sprintf("%s %s not like %v", p.str, name, val)
}

func (p *sqlProcess) OnContains(name string, value any, rValue Value) {
	if s, ok := rValue.(*StringValue); ok {
		p.str = fmt.Sprintf("%s %s like '%%%v%%'", p.str, name, s.Value)
	} else {
		panic("invalid rsql type in contains ")
	}
}

func (p *sqlProcess) OnNotContains(name string, value any, rValue Value) {
	if s, ok := rValue.(*StringValue); ok {
		p.str = fmt.Sprintf("%s %s not like '%%%v%%'", p.str, name, s.Value)
	} else {
		panic("invalid rsql type in OnNotContains ")
	}
}

func (p *sqlProcess) OnGreaterThan(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s>%v", p.str, name, val)
}

func (p *sqlProcess) OnGreaterThanOrEquals(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s>=%v", p.str, name, val)
}

func (p *sqlProcess) OnLessThan(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s<%v", p.str, name, val)
}

func (p *sqlProcess) OnLessThanOrEquals(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s <= %v", p.str, name, val)
}

func (p *sqlProcess) OnIn(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s in %s", p.str, name, val)
}

func (p *sqlProcess) OnNotIn(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s not in %v", p.str, name, val)
}

func (p *sqlProcess) OnAndItem() {
	p.str = fmt.Sprintf("%s and ", p.str)
}

func (p *sqlProcess) OnAndStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}

func (p *sqlProcess) OnAndEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}

func (p *sqlProcess) OnOrItem() {
	p.str = fmt.Sprintf("%s or ", p.str)
}

func (p *sqlProcess) OnOrStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}

func (p *sqlProcess) OnOrEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}

func (p *sqlProcess) GetSQL() string {
	return p.str
}

func (p *sqlProcess) OnIsNull(name string, value any, rValue Value) {
	p.str = fmt.Sprintf("%s %s IS NULL", p.str, name)
}

func (p *sqlProcess) OnNotIsNull(name string, value any, rValue Value) {
	p.str = fmt.Sprintf("%s %s IS NOT NULL", p.str, name)
}

func (p *sqlProcess) OnStart(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s like '%s%%'", p.str, name, val)
}

func (p *sqlProcess) OnEnd(name string, value any, rValue Value) {
	val := p.getValue(rValue)
	p.str = fmt.Sprintf("%s %s like '%%%s'", p.str, name, val)
}

func (p *sqlProcess) GetFilter() any {
	return p.GetSQL()
}

func (p *sqlProcess) getValue(value Value) any {
	var v any
	var err error
	switch value.(type) {
	case *StringValue:
		sv, _ := value.(*StringValue)
		v = "'" + sv.Value + "'"
	case *IntegerValue:
		sv, _ := value.(*IntegerValue)
		v = sv.Value
	case *DateValue:
		sv, _ := value.(*DateValue)
		v, err = time.Parse(dateLayout, sv.Value)
	case *DoubleValue:
		sv, _ := value.(*DoubleValue)
		v = sv.Value
	case *DateTimeValue:
		sv, _ := value.(*DateTimeValue)
		v, err = time.Parse(dateTimeLayout, sv.Value)
	case *BooleanValue:
		sv, _ := value.(*BooleanValue)
		v = sv.Value
	case *ListValue:
		sv, _ := value.(*ListValue)
		v = getValueList(sv)
	case *FuncValue:
		sv, _ := value.(*FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}

func (p *sqlProcess) getLikeValue(value Value) string {
	var s = ""
	if strVal, ok := value.(*StringValue); ok {
		s = strVal.Value
	} else {
		s = fmt.Sprintf("%v", value)
	}
	s = strings.Replace(s, "'", "''", -1)
	s = strings.Replace(s, "*", "%", -1)
	return "'" + s + "'"
}
