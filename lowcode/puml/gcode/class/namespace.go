package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/tags"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type Namespace struct {
	NameElement
	Elements    map[string]Element    `json:"elements,omitempty"`
	Namespaces  map[string]*Namespace `json:"namespaces,omitempty"`
	Connections []*Connection         `json:"connections,omitempty"`
	Parent      *Namespace            `json:"-"`
}

func NewNamespace(line string, parent *Namespace, notes []*Comment) *Namespace {
	ns := &Namespace{
		Elements:    make(map[string]Element),
		Connections: make([]*Connection, 0),
		Namespaces:  make(map[string]*Namespace),
		Parent:      parent,
	}
	ns.InitName(utils.ParseNamespaceName(line))
	namespace := ""
	if parent != nil {
		namespace = parent.GetNamespaceName() + "." + ns.Name
	}
	ns.InitBase(ns, line, "namespace", namespace, notes)
	return ns
}

func (n *Namespace) Parse(ctx context.Context, reader ParseReader) (resNotes []*Comment, err error) {
	notes := make([]*Comment, 0)
	var note *Note
	var element Element
	for reader.Scan() {
		line := reader.ReadLine()
		if line == "}" {
			return
		} else if strings.Contains(line, "'") {
			note := NewComment(ctx, line)
			notes = append(notes, note)
			continue
		} else if line == "" {
			continue
		} else if strings.HasPrefix(line, "note") {
			notes, element, err = n.ParseLine(ctx, line, notes, reader)
			if element != nil {
				if n, ok := element.(*Note); ok {
					note = n
				}
			}
		} else if utils.StringIsEndNote(line) {

		} else {
			notes, element, err = n.ParseLine(ctx, line, notes, reader)
			if element != nil {
				note = nil
			} else if element == nil && note != nil {
				note.AddText(line)
			} else {
				note = nil
			}
			if err != nil {
				return
			}
		}
	}
	err = n.CreateEnumsForTag(ctx)
	return
}

func (n *Namespace) GetPathName(rootPath string) string {
	pathName := n.getPaths(n)
	return rootPath + "/" + pathName
}

func (n *Namespace) getPaths(namespace *Namespace) string {
	name := ""
	if namespace != nil && namespace.Parent == nil {
		name = n.getPaths(namespace.Parent) + "/" + namespace.NamespaceName
	} else if namespace != nil {
		name = namespace.NamespaceName
	}
	return strings.ToLower(name)
}

func (n *Namespace) ParseLine(ctx context.Context, line string, comments []*Comment, reader ParseReader) (resComments []*Comment, element Element, err error) {
	resComments = nil
	element = nil
	if line == "}" || line == "" {
		return comments, element, nil
	} else if strings.HasPrefix(line, "namespace") {
		ns := NewNamespace(line, n, comments)
		n.Namespaces[ns.Name] = ns
		element = ns
		if _, err = ns.Parse(ctx, reader); err != nil {
			return
		}
	} else if strings.HasPrefix(line, "class") {
		class := NewClass(ctx, line, n.GetNamespaceName(), comments)
		element = class
		n.addElement(ctx, class)
		if err = class.Parse(ctx, reader); err != nil {
			return
		}
	} else if strings.HasPrefix(line, "enum") {
		enum := NewEnum(ctx, line, n.GetNamespaceName(), comments)
		element = enum
		n.addElement(ctx, enum)
		if err = enum.Parse(ctx, reader); err != nil {
			return
		}
	} else if strings.HasPrefix(line, "interface") {
		inter := NewInterface(ctx, line, n.GetNamespaceName(), comments)
		element = inter
		n.addElement(ctx, inter)
		if err = inter.Parse(ctx, reader); err != nil {
			return
		}
	} else if IsConnection(line) {
		conn := NewConnection(ctx, line, n.GetNamespaceName(), comments)
		element = conn
		n.addConnection(ctx, conn)
		return
	} else if strings.HasPrefix(line, "note") {
		note := NewNote(line)
		element = note
		n.addElement(ctx, note)
	} else {
		resComments = comments
		element = nil
	}
	return resComments, element, nil
}

func (n *Namespace) FindAggregateClass() []*Class {
	return n.FindClass(tags.TagTypeAggregate)
}

func (n *Namespace) FindEntityClass() []*Class {
	return n.FindClass(tags.TagTypeEntity)
}

func (n *Namespace) FindEventClass() []*Class {
	return n.FindClass(tags.TagTypeEvent)
}

func (n *Namespace) FindQueryClass() []*Class {
	return n.FindClass(tags.TagTypeQuery)
}

func (n *Namespace) FindRepositoryClass() []*Class {
	return n.FindClass(tags.TagTypeRepository)
}

func (n *Namespace) FindServiceClass() []*Class {
	return n.FindClass(tags.TagTypeService)
}

func (n *Namespace) FindViewClass() []*Class {
	return n.FindClass(tags.TagTypeView)
}

func (n *Namespace) FindCmdClass() []*Class {
	return n.FindClass(tags.TagTypeCmd)
}

func (n *Namespace) FindFormClass() []*Class {
	return n.FindClass(tags.TagTypeForm)
}

func (n *Namespace) FindEnumClass() []*Class {
	return n.FindClass(tags.TagTypeEnum)
}

func (n *Namespace) FindClass(tagType tags.TagType) []*Class {
	var entities []*Class
	for _, e := range n.Elements {
		if e.GetTypeName() == "class" {
			if class, ok := e.(*Class); ok {
				if list := class.FindTags(tagType); len(list) > 0 {
					entities = append(entities, class)
				}
			}
		}
	}
	return entities
}

func (n *Namespace) FindEnumTag() []*tags.TagEnum {
	var tagList []*tags.TagEnum
	for _, e := range n.Elements {
		if e.GetTypeName() == "class" || e.GetTypeName() == "interface" {
			if class, ok := e.(*Class); ok {
				list := class.FindTags(tags.TagTypeEnum)
				for _, item := range list {
					if v, ok := item.(*tags.TagEnum); ok {
						tagList = append(tagList, v)
					}
				}

				for _, field := range class.Fields {
					list := field.FindTags(tags.TagTypeEnum)
					for _, item := range list {
						if v, ok := item.(*tags.TagEnum); ok {
							tagList = append(tagList, v)
						}
					}
				}
			}
		}
	}
	return tagList
}

func (n *Namespace) CreateEnumsForTag(ctx context.Context) error {
	list := n.FindEnumTag()
	for _, enumTag := range list {
		el := NewEnumWithTag(ctx, enumTag.GetElement().GetDataType(), n, enumTag)
		n.addElement(ctx, el)
	}
	return nil
}

func (n *Namespace) addElement(ctx context.Context, node Element) {
	pctx := GetParseContext(ctx)
	if pctx != nil {
		pctx.AddNode(node)
	}
	n.Elements[node.GetName()] = node
}

func (n *Namespace) addConnection(ctx context.Context, conn ...*Connection) {
	n.Connections = append(n.Connections, conn...)
	n.registerConnection(ctx, conn...)
}

func (n *Namespace) registerNode(ctx context.Context, node Element) {
	pctx := GetParseContext(ctx)
	if pctx != nil {
		pctx.AddNode(node)
	}
}
func (n *Namespace) registerConnection(ctx context.Context, conn ...*Connection) {
	pctx := GetParseContext(ctx)
	if pctx != nil {
		pctx.AddConnection(conn...)
	}
}
