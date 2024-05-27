package class

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"strings"
)

type UMLBlock struct {
	FileName      string         `json:"fileName,omitempty"`
	AlisNode      map[string]any `json:"-"`
	Root          *Namespace     `json:"root"`
	allNode       map[string]Element
	allConnection []*Connection
}

func NewUMLBlock(fileName string) *UMLBlock {
	root := NewNamespace("namespace root", nil, nil)
	block := &UMLBlock{
		Root:     root,
		FileName: fileName,
	}
	return block
}

func (b *UMLBlock) Parse(ctx context.Context, reader ParseReader) (err error) {
	err = b.parse(ctx, reader, "@startuml", "@enduml")
	if err == nil {
		b.setConnection(ctx)
	}
	return err
}

func (b *UMLBlock) setConnection(ctx context.Context) {
	pctx := GetParseContext(ctx)
	if pctx != nil {
		//allConnections := pctx.GetConnections()
		//allNodes := pctx.GetNodes()
	}
}

func (b *UMLBlock) getNode(namespace string, nodeName string) Element {
	return nil
}

func (b *UMLBlock) parse(ctx context.Context, reader ParseReader, startTag string, endTag string) (err error) {
	comments := make([]*Comment, 0)
	var note *Note
	var element Element
	for reader.Scan() {
		line := reader.ReadLine()
		if strings.HasPrefix(line, startTag) {
			for reader.Scan() {
				line = reader.ReadLine()
				if line == "" {
					continue
				} else if strings.Contains(line, "'") {
					comments = append(comments, NewComment(ctx, line))
					continue
				} else if strings.HasPrefix(line, "}") {
					return nil
				} else if strings.HasPrefix(line, "note") {
					comments, element, err = b.Root.ParseLine(ctx, line, comments, reader)
					if n, ok := element.(*Note); ok {
						note = n
					}
				} else if utils.StringIsEndNote(line) {
					note = nil
				} else {
					var element Element
					comments, element, err = b.Root.ParseLine(ctx, line, comments, reader)
					if err != nil {
						return err
					} else if element == nil && note != nil {
						note.AddText(line)
					} else if element != nil {
						note = nil
					}
				}
			}
		}
		if strings.HasPrefix(line, endTag) {
			break
		}
	}
	err = b.Root.CreateEnumsForTag(ctx)
	return err
}
