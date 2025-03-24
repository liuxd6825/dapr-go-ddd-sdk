package builder

import (
	_ "embed"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/jsonschema/v6"
)

type Builder struct {
}

func NewBuilder() *Builder {
	return &Builder{}
}

// CreateHTML
//
//	@Description: generates HTML code for ui5-form based on JsonSchema
//	@receiver b
//	@param schema
//	@return string
//	@return error
func (b *Builder) CreateHTML(schema *jsonschema.Schema, templateType TemplateType, opts ...map[string]any) (string, error) {
	// 取得所有属性
	props := schema.GetSortProperties()

	// 将属性转换成Field类型
	fields := make([]Field, 0)
	for _, prop := range props {
		name := prop.Name()
		formField := Field{
			Name:     name,
			Title:    prop.Title,
			Type:     prop.GetType(),
			Required: schema.IsRequired(name),
		}
		fields = append(fields, formField)
	}

	// Prepare the context for the template
	context := pongo2.Context{
		"fields": fields,
		"opts":   b.getOpts(opts...),
		"schema": schema,
	}

	// 取得类型的模板对象
	tpl := templateType.GetTemplate()
	if tpl == nil {
		return "", fmt.Errorf("tpl cannot be nil")
	}

	// 渲染HTML
	output, err := tpl.Execute(context)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return output, nil
}

func (b *Builder) getOpts(opts ...map[string]any) map[string]any {
	opt := map[string]any{}
	for _, item := range opts {
		for k, v := range item {
			opt[k] = v
		}
	}
	return opt
}
