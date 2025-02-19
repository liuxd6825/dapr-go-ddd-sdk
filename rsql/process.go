package rsql

import (
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/common/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"time"
)

const (
	dateTimeLayout = "2006-01-02T15:04:05"
	dateLayout     = "2006-01-02"
)

type Process interface {
	OnAndItem()
	OnAndStart()
	OnAndEnd()
	OnOrItem()
	OnOrStart()
	OnOrEnd()
	OnEquals(name string, value interface{}, rValue Value)
	OnNotEquals(name string, value interface{}, rValue Value)
	OnLike(name string, value interface{}, rValue Value)
	OnNotLike(name string, value interface{}, rValue Value)
	OnGreaterThan(name string, value interface{}, rValue Value)
	OnGreaterThanOrEquals(name string, value interface{}, rValue Value)
	OnLessThan(name string, value interface{}, rValue Value)
	OnLessThanOrEquals(name string, value interface{}, rValue Value)
	OnIn(name string, value interface{}, rValue Value)
	OnNotIn(name string, value interface{}, rValue Value)
	OnContains(name string, value interface{}, rValue Value)
	OnNotContains(name string, value interface{}, rValue Value)
	OnIsNull(name string, value interface{}, rValue Value)
	OnNotIsNull(name string, value interface{}, rValue Value)
	OnStart(name string, value interface{}, rValue Value)
	OnEnd(name string, value interface{}, rValue Value)
	GetSQL() string
	GetFilter(tenantId string) (map[string]any, error)
}

/*
type SqlProcess struct {
	str string
}

func (p *SqlProcess) OnNotEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s != (%v)", p.str, name, value)
}

func (p *SqlProcess) OnLike(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like (%v)", p.str, name, value)
}

func (p *SqlProcess) OnNotLike(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not like %v", p.str, name, value)
}

func (p *SqlProcess) OnGreaterThan(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s>%v", p.str, name, value)
}

func (p *SqlProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s>=%v", p.str, name, value)
}

func (p *SqlProcess) OnLessThan(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s<%v", p.str, name, value)
}

func (p *SqlProcess) OnLessThanOrEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s <= %v", p.str, name, value)
}

func (p *SqlProcess) OnIn(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s in %v", p.str, name, value)
}

func (p *SqlProcess) OnNotIn(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not in %v", p.str, name, value)
}

func (p *SqlProcess) OnEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, value)
}

func (p *SqlProcess) NotEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, value)
}

func (p *SqlProcess) OnContains(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like (%v)", p.str, name, value)
}

func (p *SqlProcess) OnNotContains(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not like (%v)", p.str, name, value)
}

func (p *SqlProcess) OnStart(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like %v*", p.str, name, value)
}

func (p *SqlProcess) OnEnd(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like *%v", p.str, name, value)
}

func (p *SqlProcess) OnAndItem() {
	p.str = fmt.Sprintf("%s and ", p.str)
}

func (p *SqlProcess) OnAndStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}

func (p *SqlProcess) OnAndEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}
func (p *SqlProcess) OnOrItem() {
	p.str = fmt.Sprintf("%s or ", p.str)
}
func (p *SqlProcess) OnOrStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}
func (p *SqlProcess) OnOrEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}

func (p *SqlProcess) OnIsNull(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s is null", p.str, name)
}

func (p *SqlProcess) OnNotIsNull(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s is not null", p.str, name)
}

func (p *SqlProcess) Print() {
	fmt.Print(p.str)
}

func (p *SqlProcess) GetStr() string {
	return p.str
}

func NewSqlProcess() *SqlProcess {
	return &SqlProcess{str: ""}
}
func SqlParseProcess(input string) (string, error) {
	p := &SqlProcess{}
	if err := ParseProcess(input, p); err != nil {
		return "", err
	}
	return p.str, nil
} */

func ParseProcess(input string, process Process) error {
	if len(input) == 0 {
		return nil
	}
	expr, err := Parse(input)
	if err != nil {
		return errors.New(fmt.Sprintf("rsql %s expression error, %s", input, err.Error()))
	}
	err = parseProcess(expr, process)
	if err != nil {
		return errors.New(fmt.Sprintf("rsql %s parseProcess error, %s", input, err.Error()))
	}
	return nil
}

func parseProcess(expr Expression, process Process) error {
	switch expr.(type) {
	case AndExpression:
		ex, _ := expr.(AndExpression)
		process.OnAndStart()
		for i, e := range ex.Items {
			_ = parseProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnAndItem()
			}
		}
		process.OnAndEnd()
		break
	case OrExpression:
		ex, _ := expr.(OrExpression)
		process.OnOrStart()
		for i, e := range ex.Items {
			_ = parseProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnOrItem()
			}
		}
		process.OnOrEnd()
		break
	case NotEqualsComparison:
		ex, _ := expr.(NotEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotEquals(name, value, ex.Comparison.Val)
		break
	case EqualsComparison:
		ex, _ := expr.(EqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnEquals(name, value, ex.Comparison.Val)
		break
	case LikeComparison:
		ex, _ := expr.(LikeComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLike(name, value, ex.Comparison.Val)
		break
	case NotLikeComparison:
		ex, _ := expr.(NotLikeComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotLike(name, value, ex.Comparison.Val)
		break
	case GreaterThanComparison:
		ex, _ := expr.(GreaterThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThan(name, value, ex.Comparison.Val)
		break
	case GreaterThanOrEqualsComparison:
		ex, _ := expr.(GreaterThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThanOrEquals(name, value, ex.Comparison.Val)
		break
	case LessThanComparison:
		ex, _ := expr.(LessThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThan(name, value, ex.Comparison.Val)
		break
	case LessThanOrEqualsComparison:
		ex, _ := expr.(LessThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThanOrEquals(name, value, ex.Comparison.Val)
		break
	case InComparison:
		ex, _ := expr.(InComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnIn(name, value, ex.Comparison.Val)
		break
	case NotInComparison:
		ex, _ := expr.(NotInComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotIn(name, value, ex.Comparison.Val)
		break
	case ContainsComparison:
		ex, _ := expr.(ContainsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnContains(name, value, ex.Comparison.Val)
	case NotContainsComparison:
		ex, _ := expr.(NotContainsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotContains(name, value, ex.Comparison.Val)
	case NotIsNullComparison:
		ex, _ := expr.(NotIsNullComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotIsNull(name, value, ex.Comparison.Val)
	case IsNullComparison:
		ex, _ := expr.(IsNullComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnIsNull(name, value, ex.Comparison.Val)
	case StartComparison:
		ex, _ := expr.(StartComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnStart(name, value, ex.Comparison.Val)
	case EndComparison:
		ex, _ := expr.(EndComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnEnd(name, value, ex.Comparison.Val)
	}
	return nil
}

/*
func getValue(val Value) interface{} {
	var value interface{}
	switch val.(type) {
	case IntegerValue:
		value = val.(IntegerValue).Value
		break
	case BooleanValue:
		value = val.(BooleanValue).Value
		break
	case StringValue:
		value = val.(StringValue).Value
		break
	case DateTimeValue:
		value = val.(DateTimeValue).Value
		break
	case DoubleValue:
		value = val.(DoubleValue).Value
		break
	}
	return value
}
*/

func getValue(value Value) any {
	var v any
	var err error
	switch value.(type) {
	case rsql.StringValue:
		sv, _ := value.(StringValue)
		v = sv.Value
	case rsql.IntegerValue:
		sv, _ := value.(IntegerValue)
		v = sv.Value
	case rsql.DateValue:
		sv, _ := value.(DateValue)
		v, err = time.Parse(dateLayout, sv.Value)
	case rsql.DoubleValue:
		sv, _ := value.(DoubleValue)
		v = sv.Value
	case rsql.DateTimeValue:
		sv, _ := value.(DateTimeValue)
		v, err = time.Parse(dateTimeLayout, sv.Value)
	case rsql.BooleanValue:
		sv, _ := value.(BooleanValue)
		v = sv.Value
	case rsql.ListValue:
		sv, _ := value.(ListValue)
		v = getValueList(sv)
	default:
		v = value
	}
	if err != nil {
		panic(err)
	}
	return v
}

func getValueList(listValue ListValue) []any {
	list := make([]interface{}, 0)
	for _, item := range listValue.Value {
		v := getValue(item)
		list = append(list, v)
	}
	return list
}

// AsFieldName
// @Description: 转换为mongodb规范的字段名称
// @param name
// @return string
func AsFieldName(name string) string {
	return stringutils.SnakeString(name)
}
