package rsql_neo4j

import (
	"fmt"
	rsql2 "github.com/liuxd6825/dapr-go-ddd-sdk/db/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"strings"
	"time"
)

type Process struct {
	sb       strings.Builder
	dataKey  string
	tenantId string
}

func NewProcess(tenantId string, dataKey string) rsql2.Process {
	return &Process{tenantId: tenantId, dataKey: dataKey}
}

func (p *Process) GetSQL() string {
	return p.sb.String()
}

func (p *Process) GetFilter() any {
	return p.GetSQL()
}

func (p *Process) OnFnProcess(expr rsql2.Expression, fn *rsql2.FuncValue) rsql2.Value {
	switch fn.Name {
	case "sub":
		p := NewProcess(p.tenantId, p.dataKey)
		err := rsql2.ParseProcess(fn.RSQL(), p)
		if err != nil {
			panic(err)
		}
		sql := p.GetSQL()
		field := getFieldName(fn.Args["field"])
		table := fn.Args["table"]
		sql = fmt.Sprintf("(select %s from %s where %s)", field, table, sql)
		return &rsql2.StringValue{Value: sql}
	default:
		return fn
	}
	return fn
}

func (p *Process) OnEquals(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s=%v", p.dataKey, name, val))
}

func (p *Process) OnNotEquals(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s!=%v", p.dataKey, name, val))
}

func (p *Process) OnLike(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getLikeValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s=~%v", p.dataKey, name, val))
}

func (p *Process) OnNotLike(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getLikeValue(rValue)
	p.sb.WriteString(fmt.Sprintf("not %s.%s=~%v", p.dataKey, name, val))
}

func (p *Process) OnContains(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	if s, ok := rValue.(*rsql2.StringValue); ok {
		p.sb.WriteString(fmt.Sprintf("%s.%s contains '%v'", p.dataKey, name, s.Value))
	} else {
		panic("invalid rsql type in contains ")
	}
}

func (p *Process) OnNotContains(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	if s, ok := rValue.(*rsql2.StringValue); ok {
		p.sb.WriteString(fmt.Sprintf("not %s.%s contains '%v'", p.dataKey, name, s.Value))
	} else {
		panic("invalid rsql type in OnNotContains ")
	}
}

func (p *Process) OnGreaterThan(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s>%v", p.dataKey, name, val))
}

func (p *Process) OnGreaterThanOrEquals(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s>=%v", p.dataKey, name, val))
}

func (p *Process) OnLessThan(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s<%v", p.dataKey, name, val))
}

func (p *Process) OnLessThanOrEquals(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s<=%v", p.dataKey, name, val))
}

func (p *Process) OnIn(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	// 是标准数组查询
	if vList, ok := rValue.(*rsql2.ListValue); ok {
		val := p.getInValue(vList)
		p.sb.WriteString(fmt.Sprintf("%s.%s in [%s]", p.dataKey, name, val))
	} else { // 是sub子查询
		val := p.getSubSql(rValue)
		p.sb.WriteString(fmt.Sprintf("%s.%s in %s", p.dataKey, name, val))
	}
}

func (p *Process) OnNotIn(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	// 是标准数组查询
	if vList, ok := rValue.(*rsql2.ListValue); ok {
		val := p.getInValue(vList)
		p.sb.WriteString(fmt.Sprintf("not %s.%s in [%s]", p.dataKey, name, val))
	} else { // 是sub子查询
		val := p.getSubSql(rValue)
		p.sb.WriteString(fmt.Sprintf("not %s.%s in %s", p.dataKey, name, val))
	}
}
func (p *Process) getSubSql(rValue rsql2.Value) string {
	val := p.getValue(rValue)
	if sub, ok := val.(string); ok {
		sub = strings.Trim(sub, "'")
		return sub
	}
	panic("invalid rsql type in sub sql")
}

func (p *Process) OnIsNull(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	p.sb.WriteString(fmt.Sprintf("%s.%s is null", p.dataKey, name))
}

func (p *Process) OnNotIsNull(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	p.sb.WriteString(fmt.Sprintf("not %s.%s is null", p.dataKey, name))
}

func (p *Process) OnStart(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s=~'%s.*'", p.dataKey, name, val))
}

func (p *Process) OnEnd(name string, value any, rValue rsql2.Value) {
	name = getFieldName(name)
	val := p.getValue(rValue)
	p.sb.WriteString(fmt.Sprintf("%s.%s=~'*.%s'", p.dataKey, name, val))
}

func (p *Process) OnAndItem() {
	p.sb.WriteString(" and ")
}

func (p *Process) OnOrItem() {
	p.sb.WriteString(" or ")
}
func (p *Process) OnAndStart() {
	p.sb.WriteString("(")
}

func (p *Process) OnAndEnd() {
	p.sb.WriteString(")")
}

func (p *Process) OnOrStart() {
	p.sb.WriteString("(")
}

func (p *Process) OnOrEnd() {
	p.sb.WriteString(")")
}

func (p *Process) getValue(value rsql2.Value) any {
	var v any
	var err error
	switch value.(type) {
	case *rsql2.StringValue:
		sv, _ := value.(*rsql2.StringValue)
		//v = sv.Value
		v = "'" + sv.Value + "'"
	case *rsql2.IntegerValue:
		sv, _ := value.(*rsql2.IntegerValue)
		v = sv.Value
	case *rsql2.DateValue:
		sv, _ := value.(*rsql2.DateValue)
		if date, er := time.Parse(rsql2.DateLayout, sv.Value); er != nil {
			err = er
		} else {
			v = fmt.Sprintf("'%d-%02d-%02d'", date.Year(), date.Month(), date.Day())
		}
	case *rsql2.DoubleValue:
		sv, _ := value.(*rsql2.DoubleValue)
		v = sv.Value
	case *rsql2.DateTimeValue:
		sv, _ := value.(*rsql2.DateTimeValue)
		if date, er := time.Parse(rsql2.DateTimeLayout, sv.Value); er != nil {
			err = er
		} else {
			v = fmt.Sprintf("'%d-%02d-%02d %02d:%02d:%02d'", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second())
		}
	case *rsql2.BooleanValue:
		sv, _ := value.(*rsql2.BooleanValue)
		v = sv.Value
	case *rsql2.ListValue:
		sv, _ := value.(*rsql2.ListValue)
		v = p.getValueList(sv)
	case *rsql2.FuncValue:
		sv, _ := value.(*rsql2.FuncValue)
		v = sv.Value
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}

func (p *Process) getLikeValue(value rsql2.Value) string {
	var s = ""
	if strVal, ok := value.(*rsql2.StringValue); ok {
		s = strVal.Value
	} else {
		s = fmt.Sprintf("%v", value)
	}
	s = strings.Replace(s, "'", "''", -1)
	s = strings.Replace(s, "*", ".*", -1)
	return "'" + s + "'"
}

func (p *Process) getValueList(listValue *rsql2.ListValue) []interface{} {
	list := make([]interface{}, 0)
	if listValue == nil {
		return list
	}
	for _, v := range listValue.Value {
		list = append(list, p.getValue(v))
	}
	return list
}

func (p *Process) getInValue(listValue *rsql2.ListValue) string {
	sb := &strings.Builder{}
	count := len(listValue.Value)
	for i, v := range listValue.Value {
		if s, ok := v.(*rsql2.StringValue); ok {
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

func getFieldName(fieldName string) string {
	return stringutils.FirstLowerCamelString(fieldName)
}
