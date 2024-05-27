package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Class struct {
	NameElement
	Fields      []*Field  `json:"fields,omitempty"`
	Methods     []*Method `json:"methods,omitempty"`
	GenericType []string  `json:"genericType,omitempty"` // 是泛型类
	IsInterface bool      `json:"isInterface"`
}

func NewClass(ctx context.Context, line string, namespace string, notes []*Comment) *Class {
	class := &Class{
		Fields:  make([]*Field, 0),
		Methods: make([]*Method, 0),
	}
	ns := namespace
	parseRes := ParseClassName(line, "class")
	if parseRes.NamespaceName != "" {
		ns = ns + "." + parseRes.NamespaceName
	}
	class.InitName(parseRes.Name, parseRes.Alias)
	class.InitBase(class, line, "class", ns, notes)
	class.InitDataTag(parseRes.Name)
	class.GenericType = parseRes.GenericType
	return class
}

func NewInterface(ctx context.Context, line string, namespace string, notes []*Comment) *Class {
	class := &Class{
		Fields:  make([]*Field, 0),
		Methods: make([]*Method, 0),
	}
	ns := namespace
	parseRes := ParseClassName(line, "interface")
	if parseRes.NamespaceName != "" {
		ns = ns + "." + parseRes.NamespaceName
	}
	class.InitName(parseRes.Name, parseRes.Alias)
	class.InitBase(class, line, "interface", ns, notes)
	class.GenericType = parseRes.GenericType
	class.IsInterface = true
	return class
}

func (c *Class) Parse(ctx context.Context, reader ParseReader) error {

	if !strings.HasSuffix(c.Line, "{") {
		return nil
	}
	/*
		{field} A field (despite parentheses)
		{method} Some method
	*/
	var notes []*Comment
	for reader.Scan() {
		line := reader.ReadLine()
		if line == "}" {
			return nil
		} else if strings.HasPrefix(line, "==") {
			continue
		} else if strings.HasPrefix(line, "..") {
			continue
		} else if strings.HasPrefix(line, "__") {
			continue
		} else if strings.HasPrefix(line, "--") {
			continue
		} else if strings.HasPrefix(line, "'") {
			notes = append(notes, NewComment(ctx, line))
		} else if utils.StringContains(line, "(", "{method}") {
			c.AddMethod(NewMethod(ctx, c, line, notes).Init(line))
			notes = nil
		} else if len(line) > 0 {
			c.AddField(NewField(ctx, c, line, c.NamespaceName, notes))
			notes = nil
		}
	}
	return nil
}

func (c *Class) AddField(member *Field) {
	c.Fields = append(c.Fields, member)
}

func (c *Class) AddMethod(member *Method) {
	c.Methods = append(c.Methods, member)
}

func (c *Class) LenFields() int {
	return len(c.Fields)
}
func (c *Class) LenMethods() int {
	return len(c.Methods)
}

type ResultParseClassName struct {
	Name          string   `json:"name"`
	Alias         string   `json:"alias"`
	GenericType   []string `json:"genericType"`
	Extends       []string `json:"extends"`
	Implements    []string `json:"implements"`
	NamespaceName string   `json:"namespaceName"`
}

const extends = "extends"
const implements = "implements"

func ParseClassName(line string, tagName string) *ResultParseClassName {
	name, alias, namespaceName, list := utils.ParseName(line, tagName)
	res := &ResultParseClassName{
		Name:          name,
		Alias:         alias,
		NamespaceName: namespaceName,
	}
	res.setExtends(line)
	res.setImplements(line)
	res.setGenericType(list)
	return res
}

func (r *ResultParseClassName) setExtends(line string) {
	if strings.Contains(line, extends) {
		idx := strings.Index(line, extends)
		s := line[idx+7:]
		s = strings.ReplaceAll(s, ",", " ")
		list := strings.Split(s, " ")
		for _, item := range list {
			item = strings.ReplaceAll(item, "{", "")
			if item == implements {
				break
			}
			if item != "" {
				r.Extends = append(r.Extends, item)
			}
		}
	}
}

func (r *ResultParseClassName) setImplements(line string) {
	if strings.Contains(line, implements) {
		idx := strings.Index(line, implements)
		s := line[idx+10:]
		s = strings.ReplaceAll(s, ",", " ")
		list := strings.Split(s, " ")
		for _, item := range list {
			item = strings.ReplaceAll(item, "{", "")
			if item == extends {
				break
			}
			if item != "" {
				r.Implements = append(r.Implements, item)
			}
		}
	}
}

func (r *ResultParseClassName) setGenericType(list []string) {
	isGenericType := false
	for _, item := range list {
		// 解析 <<User>>
		if strings.HasPrefix(item, "<<") && strings.HasSuffix(item, ">>") {
			str := strings.Trim(item[2:len(item)-2], " ")
			str = strings.ReplaceAll(item, ",", " ")
			r.GenericType = strings.Split(str, " ")
			break
		}
		// 解析 << User >>
		if strings.HasPrefix(item, "<<") {
			isGenericType = true
			s := strings.ReplaceAll(item[2:], ",", "")
			if len(s) > 0 {
				r.AddGenericType(s)
			}
		} else if strings.HasSuffix(item, ">>") {
			isGenericType = false
			s := strings.ReplaceAll(item[:len(item)-2], ",", "")
			if len(s) > 0 {
				r.AddGenericType(s)
			}
			break
		} else if isGenericType {
			s := strings.ReplaceAll(item, ",", "")
			if len(s) > 0 {
				r.AddGenericType(s)
			}
		}
	}
}

func (r *ResultParseClassName) AddGenericType(typeName string) {
	s := strings.Trim(typeName, " ")
	if s != "" && s != "," && s != ";" && s != extends && s != implements {
		r.GenericType = append(r.GenericType, s)
	}
}
