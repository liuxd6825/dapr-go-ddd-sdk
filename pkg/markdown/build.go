package markdown

import (
	"fmt"
	"github.com/liuxd6825/jsonschema/v6"
	"strings"
)

type Builder struct {
}

func NewBuilder() *Builder {
	return &Builder{}
}

// Build 根据JSON Schema生成Markdown文档
func (b *Builder) Build(data any, sch *jsonschema.Schema) (sb *strings.Builder, err error) {
	sb = &strings.Builder{}
	if m, ok := data.(map[string]any); ok {
		err = b.AddObject(m, sch, sb, 1)
	} else if list, ok := data.([]any); ok {
		err = b.AddList(list, sch, sb, 1)
	} else if list, ok := data.([]map[string]any); ok {
		err = b.AddList(list, sch, sb, 1)
	}
	return sb, nil
}

// List 根据JSON Schema生成Markdown文档
func (b *Builder) AddList(data any, sch *jsonschema.Schema, builder *strings.Builder, level int) error {
	// 写入标题
	builder.WriteString(fmt.Sprintf("%s %s\n\n", b.getLevel(level), sch.Title))

	// 写入描述
	if desc := sch.Description; desc != "" {
		builder.WriteString(fmt.Sprintf("%s\n\n", desc))
	}
	builder.WriteString(fmt.Sprintf("%s %s\n\n", b.getLevel(level+1), sch.Title))

	props := sch.GetSortProperties()
	var fields []string
	// 处理属性
	for _, prop := range props {
		fields = append(fields, prop.Title)
	}

	builder.WriteString("| " + strings.Join(fields, "| ") + " |\n")
	builder.WriteString(strings.Repeat("|-------", len(fields)) + " |\n")
	if list, ok := data.([]map[string]any); ok {
		for _, m := range list {
			builder.WriteString("|")
			for _, prop := range props {
				val := b.getMapValue(prop, m)
				builder.WriteString(fmt.Sprintf("%v|", val))
			}
			builder.WriteString("\n")
		}
	} else if list, ok := data.([]any); ok {
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				builder.WriteString("|")
				for _, prop := range props {
					val := b.getMapValue(prop, m)
					builder.WriteString(fmt.Sprintf("%v|", val))
				}
				builder.WriteString("\n")
			}
		}
	}

	return nil
}

func (b *Builder) getMapValue(prop *jsonschema.Schema, m map[string]any) any {
	val := m[prop.Name()]
	if val == nil {
		val = ""
	}
	return val
}

// Object 根据JSON Schema生成Markdown文档
func (b *Builder) AddObject(data any, sch *jsonschema.Schema, builder *strings.Builder, level int) error {
	// 写入标题
	builder.WriteString(fmt.Sprintf("%s %s\n\n", b.getLevel(level), sch.Title))

	// 写入描述
	if desc := sch.Description; desc != "" {
		builder.WriteString(fmt.Sprintf("%s\n\n", desc))
	}
	mapData, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("data should be an object")
	}
	props := sch.GetSortProperties()
	var listProps []*jsonschema.Schema
	var objectProps []*jsonschema.Schema
	// 输出属性
	for _, prop := range props {
		if prop.Types == nil {
			continue
		}
		if prop.Types.Contains(jsonschema.JsonType_ArrayType) {
			listProps = append(listProps, prop)
		} else if prop.Types.Contains(jsonschema.JsonType_ObjectType) {
			objectProps = append(objectProps, prop)
		} else {
			title := prop.Title
			val := b.getMapValue(prop, mapData)
			builder.WriteString(fmt.Sprintf("- **%s**: %v\n", title, val))
		}
	}

	// 输出对象
	for _, prop := range objectProps {
		val := b.getMapValue(prop, mapData)
		if list, ok := val.([]any); ok {
			if err := b.AddList(list, prop, builder, level+1); err != nil {
				return err
			}
		}
	}

	// 输出明细表
	for _, prop := range listProps {
		val := b.getMapValue(prop, mapData)
		err := b.AddObject(val, prop, builder, level+1)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) getLevel(count int) string {
	return strings.Repeat("#", count)
}
