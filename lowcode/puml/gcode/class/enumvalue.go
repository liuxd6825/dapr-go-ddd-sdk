package class

import "context"

type EnumValue struct {
	NameElement
	Line string `json:"-"`
}

func NewEnumValue(ctx context.Context, line string, namespaceName string, notes []*Comment) *EnumValue {
	member := &EnumValue{}
	member.Name = line
	member.Line = line
	member.InitBase(member, line, "EnumValue", namespaceName, notes)
	return member
}
