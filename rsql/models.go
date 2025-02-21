package rsql

import (
	"fmt"
)

type Identifier struct {
	Val string
}

type Value interface{ ValueName() string }

type StringValue struct {
	Value string
}

func (s *StringValue) ValueName() string { return "string" }

func (s *StringValue) String() string { return s.Value }

type BooleanValue struct{ Value bool }

func (v *BooleanValue) ValueName() string { return "bool" }
func (v *BooleanValue) String() string {
	if v.Value {
		return "true"
	}
	return "false"
}

type DateValue struct{ Value string }

func (v *DateValue) ValueName() string { return "date" }
func (v *DateValue) String() string {
	return v.Value
}

type DateTimeValue struct{ Value string }

func (v *DateTimeValue) ValueName() string { return "datetime" }
func (v *DateTimeValue) String() string {
	return v.Value
}

type IntegerValue struct{ Value int64 }

func (v *IntegerValue) ValueName() string { return "int" }
func (v *IntegerValue) String() string {
	return fmt.Sprintf("%d", v.Value)
}

type DoubleValue struct{ Value float64 }

func (v *DoubleValue) ValueName() string { return "double" }
func (v *DoubleValue) String() string {
	return fmt.Sprintf("%f", v.Value)
}

type ListValue struct{ Value []Value }

func (v *ListValue) ValueName() string { return "list" }
func (v *ListValue) String() string {
	return fmt.Sprintf("%v", v.Value)
}

type Expression interface{ ExpressionName() string }

type OrExpression struct{ Items []Expression }

func (*OrExpression) ExpressionName() string { return "Or" }

type AndExpression struct{ Items []Expression }

func (*AndExpression) ExpressionName() string { return "And" }

type Comparison struct {
	Identifier Identifier
	Val        Value
}

func (c *Comparison) ExpressionName() string { return "Comparison" }

func (c *Comparison) Value() Value {
	return c.Val
}
func (c *Comparison) SetValue(val Value) {
	c.Val = val
}

func (c *Comparison) GetIdentifier() Identifier {
	return c.Identifier
}
func (c *Comparison) SetIdentifier(val Identifier) {
	c.Identifier = val
}

type GetIdentifier interface {
	GetIdentifier() Identifier
}

type SetIdentifier interface {
	SetIdentifier(val Identifier)
}

type EqualsComparison struct{ Comparison }

func (c *EqualsComparison) ExpressionName() string { return "==" }

type NotEqualsComparison struct{ Comparison }

func (c *NotEqualsComparison) ExpressionName() string { return "!=" }

type LikeComparison struct{ Comparison }

func (c *LikeComparison) ExpressionName() string { return "~=" }

type NotLikeComparison struct{ Comparison }

func (c *NotLikeComparison) ExpressionName() string { return "!~=" }

type GreaterThanComparison struct{ Comparison }

func (c *GreaterThanComparison) ExpressionName() string { return ">" }

type GreaterThanOrEqualsComparison struct{ Comparison }

func (c *GreaterThanOrEqualsComparison) ExpressionName() string { return ">=" }

type LessThanComparison struct{ Comparison }

func (c *LessThanComparison) ExpressionName() string { return "<" }

type LessThanOrEqualsComparison struct{ Comparison }

func (c *LessThanOrEqualsComparison) ExpressionName() string { return "<=" }

type InComparison struct{ Comparison }

func (c *InComparison) ExpressionName() string { return "=in=" }

type NotInComparison struct{ Comparison }

func (c *NotInComparison) ExpressionName() string { return "=out=" }

type ContainsComparison struct{ Comparison }

func (c *ContainsComparison) ExpressionName() string { return "=contains=" }

type NotContainsComparison struct{ Comparison }

func (c *NotContainsComparison) ExpressionName() string { return "=!contains=" }

type IsNullComparison struct{ Comparison }

func (c *IsNullComparison) ExpressionName() string { return "=null=" }

type NotIsNullComparison struct{ Comparison }

func (c *NotIsNullComparison) ExpressionName() string { return "=!null=" }

type StartComparison struct{ Comparison }

func (c *StartComparison) ExpressionName() string { return "=start=" }

type EndComparison struct{ Comparison }

func (c *EndComparison) ExpressionName() string { return "=end=" }

type FuncComparison struct{ Comparison }

func (c *FuncComparison) ExpressionName() string { return "=func=" }
