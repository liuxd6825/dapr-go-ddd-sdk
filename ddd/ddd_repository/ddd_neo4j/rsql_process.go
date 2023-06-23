package ddd_neo4j

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
)

type RsqlProcess interface {
	OnAndItem()
	OnAndStart()
	OnAndEnd()
	OnOrItem()
	OnOrStart()
	OnOrEnd()
	OnEquals(name string, value interface{}, rValue rsql.Value)
	OnNotEquals(name string, value interface{}, rValue rsql.Value)
	OnLike(name string, value interface{}, rValue rsql.Value)
	OnNotLike(name string, value interface{}, rValue rsql.Value)
	OnGreaterThan(name string, value interface{}, rValue rsql.Value)
	OnGreaterThanOrEquals(name string, value interface{}, rValue rsql.Value)
	OnLessThan(name string, value interface{}, rValue rsql.Value)
	OnLessThanOrEquals(name string, value interface{}, rValue rsql.Value)
	OnIn(name string, value interface{}, rValue rsql.Value)
	OnNotIn(name string, value interface{}, rValue rsql.Value)
	OnContains(name string, value interface{}, rValue rsql.Value)
	OnNotContains(name string, value interface{}, rValue rsql.Value)

	GetSqlWhere(tenantId string) interface{}
	GetFilter(tenantId string) interface{}
}

type rsqlProcess struct {
	str     string
	dataKey string
}

func NewRSqlProcess(dataKey string) RsqlProcess {
	return &rsqlProcess{dataKey: dataKey}
}

func (p *rsqlProcess) OnNotEquals(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s != (%v)", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnLike(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s like (%v)", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnNotLike(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s not like %v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnGreaterThan(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s>%v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnGreaterThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s>=%v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnLessThan(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s<%v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnLessThanOrEquals(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s <= %v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnIn(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s in %v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnNotIn(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s not in %v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) OnEquals(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s=%v", p.str, p.dataKey, name, value)
}

func (p *rsqlProcess) NotEquals(name string, value interface{}, rValue rsql.Value) {
	p.str = fmt.Sprintf("%s %s.%s=%v", p.str, p.dataKey, name, value)
}

func (r *rsqlProcess) OnContains(name string, value interface{}, rValue rsql.Value) {
	//TODO implement me
	panic("implement me")
}

func (r *rsqlProcess) OnNotContains(name string, value interface{}, rValue rsql.Value) {
	//TODO implement me
	panic("implement me")
}

func (p *rsqlProcess) OnAndItem() {
	p.str = fmt.Sprintf("%s and ", p.str)
}
func (p *rsqlProcess) OnAndStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}
func (p *rsqlProcess) OnAndEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}
func (p *rsqlProcess) OnOrItem() {
	p.str = fmt.Sprintf("%s or ", p.str)
}
func (p *rsqlProcess) OnOrStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}
func (p *rsqlProcess) OnOrEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}

func (p *rsqlProcess) GetFilter(tenantId string) interface{} {
	return p.str
}

func (p *rsqlProcess) GetSqlWhere(tenantId string) interface{} {
	return p.str
}

func (p *rsqlProcess) Print() {
	fmt.Print(p.str)
}

func getNeo4jWhere(tenantId string, dataKey string, filter string) (string, error) {
	p := NewRSqlProcess(dataKey)

	if err := ParseProcess(filter, p); err != nil {
		return "", err
	}
	where := p.GetSqlWhere(tenantId).(string)
	if len(where) > 0 {
		where = "WHERE " + where
	}
	return where, nil
}

func ParseProcess(input string, process RsqlProcess) error {
	if len(input) == 0 {
		return nil
	}
	expr, err := rsql.Parse(input)
	if err != nil {
		return errors.New(fmt.Sprintf("rsql %s expression error, %s", input, err.Error()))
	}
	err = parseRSqlProcess(expr, process)
	if err != nil {
		return errors.New(fmt.Sprintf("rsql %s parseRSqlProcess error, %s", input, err.Error()))
	}
	return nil
}

func parseRSqlProcess(expr rsql.Expression, process rsql.Process) error {
	switch expr.(type) {
	case rsql.AndExpression:
		ex, _ := expr.(rsql.AndExpression)
		process.OnAndStart()
		for i, e := range ex.Items {
			_ = parseRSqlProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnAndItem()
			}
		}
		process.OnAndEnd()
		break
	case rsql.OrExpression:
		ex, _ := expr.(rsql.OrExpression)
		process.OnOrStart()
		for i, e := range ex.Items {
			_ = parseRSqlProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnOrItem()
			}
		}
		process.OnOrEnd()
		break
	case rsql.NotEqualsComparison:
		ex, _ := expr.(rsql.NotEqualsComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnNotEquals(name, value, ex.Comparison.Val)
		break
	case rsql.EqualsComparison:
		ex, _ := expr.(rsql.EqualsComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnEquals(name, value, ex.Comparison.Val)
		break
	case rsql.LikeComparison:
		ex, _ := expr.(rsql.LikeComparison)
		name := getFieldName(ex.Comparison.Identifier.Val)
		value := getValue(ex.Comparison.Val)
		process.OnLike(name, value, ex.Comparison.Val)
		break
	case rsql.NotLikeComparison:
		ex, _ := expr.(rsql.NotLikeComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnNotLike(name, value, ex.Comparison.Val)
		break
	case rsql.GreaterThanComparison:
		ex, _ := expr.(rsql.GreaterThanComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThan(name, value, ex.Comparison.Val)
		break
	case rsql.GreaterThanOrEqualsComparison:
		ex, _ := expr.(rsql.GreaterThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThanOrEquals(name, value, ex.Comparison.Val)
		break
	case rsql.LessThanComparison:
		ex, _ := expr.(rsql.LessThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThan(name, value, ex.Comparison.Val)
		break
	case rsql.LessThanOrEqualsComparison:
		ex, _ := expr.(rsql.LessThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnLessThanOrEquals(name, value, ex.Comparison.Val)
		break
	case rsql.InComparison:
		ex, _ := expr.(rsql.InComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnIn(name, value, ex.Comparison.Val)
		break
	case rsql.NotInComparison:
		ex, _ := expr.(rsql.NotInComparison)
		name := ex.Comparison.Identifier.Val
		name = getFieldName(name)
		value := getValue(ex.Comparison.Val)
		process.OnNotIn(name, value, ex.Comparison.Val)
		break
	}
	return nil
}

func getValue(val rsql.Value) interface{} {
	var value interface{}
	switch val.(type) {
	case rsql.IntegerValue:
		value = val.(rsql.IntegerValue).Value
		break
	case rsql.BooleanValue:
		value = val.(rsql.BooleanValue).Value
		break
	case rsql.StringValue:
		value = fmt.Sprintf(`'%v'`, val.(rsql.StringValue).Value)
		break
	case rsql.DateTimeValue:
		value = fmt.Sprintf(`'%v'`, val.(rsql.DateTimeValue).Value)
		break
	case rsql.DoubleValue:
		value = val.(rsql.DoubleValue).Value
		break
	}
	return value
}

func getFieldName(name string) string {
	return name
}
