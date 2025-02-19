package rsql

import (
	"fmt"
	"strings"
)

type sqlProcess struct {
	str string
}

func NewSqlProcess() Process {
	return &sqlProcess{}
}

func (p *sqlProcess) OnNotEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s != (%v)", p.str, name, value)
}

func (p *sqlProcess) OnLike(name string, value interface{}, rValue Value) {
	val := getLikeValue(value)
	p.str = fmt.Sprintf("%s %s like '%v'", p.str, name, val)
}

func getLikeValue(value interface{}) string {
	var val string
	if strVal, ok := value.(StringValue); ok {
		val = strVal.Value
	} else {
		val = fmt.Sprintf("%v", value)
	}
	val = strings.Replace(val, "'", "''", -1)
	val = strings.Replace(val, "*", "%", -1)
	return val
}

func (p *sqlProcess) OnNotLike(name string, value interface{}, rValue Value) {
	val := getLikeValue(value)
	p.str = fmt.Sprintf("%s %s not like '%v'", p.str, name, val)
}

func (p *sqlProcess) OnGreaterThan(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s>%v", p.str, name, value)
}

func (p *sqlProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s>=%v", p.str, name, value)
}

func (p *sqlProcess) OnLessThan(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s<%v", p.str, name, value)
}

func (p *sqlProcess) OnLessThanOrEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s <= %v", p.str, name, value)
}

func (p *sqlProcess) OnIn(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s in %v", p.str, name, value)
}

func (p *sqlProcess) OnNotIn(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not in %v", p.str, name, value)
}

func (p *sqlProcess) OnEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, value)
}

func (p *sqlProcess) NotEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, value)
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

func (p *sqlProcess) OnContains(name string, value interface{}, rValue Value) {
	val := getLikeValue(value)
	p.str = fmt.Sprintf("%s %s like '%%%v%%'", p.str, name, val)
}

func (p *sqlProcess) OnNotContains(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not like '%%%v%%'", p.str, name, value)
}

func (p *sqlProcess) OnIsNull(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s IS NULL", p.str, name, value)
}

func (p *sqlProcess) OnNotIsNull(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s IS NOT NULL", p.str, name, value)
}

func (p *sqlProcess) OnStart(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like '%s%%'", p.str, name, value)
}

func (p *sqlProcess) OnEnd(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like '%%%s'", p.str, name, value)
}

func (p *sqlProcess) GetFilter(tenantId string) (map[string]any, error) {
	return map[string]any{}, nil
}
