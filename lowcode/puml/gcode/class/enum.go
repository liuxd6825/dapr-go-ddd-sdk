package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/tags"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Enum struct {
	NameElement
	Members []*EnumValue `json:"members,omitempty"`
}

func NewEnum(ctx context.Context, line string, namespace string, notes []*Comment) *Enum {
	enum := &Enum{
		Members: make([]*EnumValue, 0),
	}
	name, alias, ns := utils.ParseEnumName(line)
	if namespace != "" && ns != "" {
		ns = "." + ns
	}
	enum.InitBase(enum, line, "enum", namespace+ns, notes)
	enum.InitName(name, alias)
	enum.InitDataTag(name)
	return enum
}

func NewEnumWithTag(ctx context.Context, name string, namespace *Namespace, tag *tags.TagEnum) *Enum {
	enum := &Enum{
		Members: make([]*EnumValue, 0),
	}
	enum.InitBase(enum, "", "enum", namespace.NamespaceName, nil)
	enum.InitName(name, "")
	enum.Name = tag.Element.GetDataType()
	enum.NamespaceName = namespace.NamespaceName
	for _, item := range tag.Items {
		enumValue := NewEnumValue(ctx, item.Name, namespace.NamespaceName, nil)
		enum.AddComment(item.Title)
		enum.AddEnumValue(enumValue)
	}
	return enum
}

func (e *Enum) Parse(ctx context.Context, reader ParseReader) error {
	var notes []*Comment
	for reader.Scan() {
		line := reader.ReadLine()
		if line == "}" {
			return nil
		} else if strings.HasPrefix(line, "'") {
			notes = append(notes, NewComment(ctx, line))
		} else {
			e.AddEnumValue(NewEnumValue(ctx, line, e.GetNamespaceName(), notes))
			notes = nil
		}
	}
	return nil
}

func (e *Enum) AddEnumValue(member *EnumValue) {
	e.Members = append(e.Members, member)
}
