package code_generator

import (
	"text/template"

	"github.com/liuxd6825/jsonschema/v6"
)

const ModelTemplate = `package {{.PackageName}}

import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)

type {{.StructName}} struct {
{{- if .HasBase}}
	xbase.BaseModel ` + "`bson:\",inline\"`" + `
{{- end}}
{{- range .Fields}}
	{{.Name}} {{.Type}} ` + "`{{.Tag}}`" + `  
{{- end}}
}

func New{{.StructName}}() *{{.StructName}} {
	return &{{.StructName}}{}
}
`

var modelTpl *template.Template

type ModelOption struct {
}

// GeneratorModel 处理 JsonSchema 并生成代码
func (g *Generator) GeneratorModel(schema *jsonschema.Schema, opts ...*Options) ([]byte, error) {
	opt := NewOptions(opts...)
	if opt.Template == nil {
		opt.Template = GetModelTpl()
	}
	if opt.PackageName == "" {
		opt.PackageName = "model"
	}
	return g.Generate(opt.Template, schema, opt.PackageName, opt.Values)
}

func GetModelTpl() *template.Template {
	if modelTpl == nil {
		modelTpl = template.Must(template.New("model").Parse(ModelTemplate))
	}
	return modelTpl
}
