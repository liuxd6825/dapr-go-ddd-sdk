package rsql

import (
	"errors"
	"fmt"
)

const (
	DateTimeLayout = "2006-01-02T15:04:05"
	DateLayout     = "2006-01-02"
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

	OnFnProcess(expr Expression, fn *FuncValue) Value
	GetSQL() string
	GetFilter() any
}

func ParseProcess(input string, process Process) error {
	if len(input) == 0 {
		return nil
	}
	expr, err := Parse(input)
	if err != nil {
		return errors.New(fmt.Sprintf("%s expression error, %s", input, err.Error()))
	}
	err = parseProcess(expr, process)
	if err != nil {
		return errors.New(fmt.Sprintf("%s parseProcess error, %s", input, err.Error()))
	}
	return nil
}

type ValueComparison interface {
	Value() Value
	SetValue(value Value)
}

func parseProcess(expr Expression, process Process) error {
	if v, ok := expr.(ValueComparison); ok {
		val := v.Value()
		if fn, ok := val.(*FuncValue); ok {
			val = process.OnFnProcess(expr, fn)
			if val == nil {
				return nil
			}
			v.SetValue(val)
		}
	}
	switch expr.(type) {
	case *FuncComparison:
		ex, _ := expr.(*FuncComparison)
		fmt.Println(ex.Val)
		break
	case *AndExpression:
		ex, _ := expr.(*AndExpression)
		process.OnAndStart()
		for i, e := range ex.Items {
			_ = parseProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnAndItem()
			}
		}
		process.OnAndEnd()
		break
	case *OrExpression:
		ex, _ := expr.(*OrExpression)
		process.OnOrStart()
		for i, e := range ex.Items {
			_ = parseProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnOrItem()
			}
		}
		process.OnOrEnd()
		break
	case *NotEqualsComparison:
		ex, _ := expr.(*NotEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotEquals(name, value, ex.Comparison.Val)
		break
	case *EqualsComparison:
		ex, _ := expr.(*EqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnEquals(name, value, ex.Comparison.Val)
		break
	case *LikeComparison:
		ex, _ := expr.(*LikeComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLike(name, value, ex.Comparison.Val)
		break
	case *NotLikeComparison:
		ex, _ := expr.(*NotLikeComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotLike(name, value, ex.Comparison.Val)
		break
	case *GreaterThanComparison:
		ex, _ := expr.(*GreaterThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThan(name, value, ex.Comparison.Val)
		break
	case *GreaterThanOrEqualsComparison:
		ex, _ := expr.(*GreaterThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThanOrEquals(name, value, ex.Comparison.Val)
		break
	case *LessThanComparison:
		ex, _ := expr.(*LessThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThan(name, value, ex.Comparison.Val)
		break
	case *LessThanOrEqualsComparison:
		ex, _ := expr.(*LessThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThanOrEquals(name, value, ex.Comparison.Val)
		break
	case *InComparison:
		ex, _ := expr.(*InComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnIn(name, value, ex.Comparison.Val)
		break
	case *NotInComparison:
		ex, _ := expr.(*NotInComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotIn(name, value, ex.Comparison.Val)
		break
	case *ContainsComparison:
		ex, _ := expr.(*ContainsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnContains(name, value, ex.Comparison.Val)
	case *NotContainsComparison:
		ex, _ := expr.(*NotContainsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotContains(name, value, ex.Comparison.Val)
	case *NotIsNullComparison:
		ex, _ := expr.(*NotIsNullComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotIsNull(name, value, ex.Comparison.Val)
	case *IsNullComparison:
		ex, _ := expr.(*IsNullComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnIsNull(name, value, ex.Comparison.Val)
	case *StartComparison:
		ex, _ := expr.(*StartComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnStart(name, value, ex.Comparison.Val)
	case *EndComparison:
		ex, _ := expr.(*EndComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnEnd(name, value, ex.Comparison.Val)
	}
	return nil
}
