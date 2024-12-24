package form

import (
	_ "embed"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/jsonschema/v6"
)

type Field struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type Schema struct {
	Fields []Field `json:"fields"`
}

// Embed the HTML template file
//
//go:embed form.tpl.html
var ui5FormTemplate string

type Builder struct {
	tmpl *pongo2.Template
}

func NewBuilder() *Builder {
	tmpl, err := pongo2.FromString(ui5FormTemplate)
	if err != nil {
		panic(err)
	}
	return &Builder{tmpl: tmpl}
}

// CreateHTML generates HTML code for ui5-form based on JsonSchema
func (b *Builder) CreateHTML(schema *jsonschema.Schema) (string, error) {

	props := schema.GetAllProperties()
	formSchema := Schema{Fields: make([]Field, 0)}
	for key, prop := range props {
		formField := Field{
			Name:     key,
			Title:    prop.Title,
			Type:     prop.GetType(),
			Required: schema.IsRequired(key),
		}
		formSchema.Fields = append(formSchema.Fields, formField)
	}

	// Prepare the context for the template
	context := pongo2.Context{
		"fields":     formSchema.Fields,
		"formSchema": formSchema,
	}

	// Render the template with the provided schema
	output, err := b.tmpl.Execute(context)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return output, nil
}
