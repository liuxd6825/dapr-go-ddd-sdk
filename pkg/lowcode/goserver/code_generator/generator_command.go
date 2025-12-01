package code_generator

import (
	"text/template"

	"github.com/liuxd6825/jsonschema/v6"
)

const CommandTemplate = `
package {{.PackageName}}

import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)

type {{.StructName}} struct {
{{- if .HasBase}}
	xbase.BaseModel bson:\",inline\"   
{{- end}}
{{- range .Fields}}
	{{.Name}} {{.Type}}   {{.Tag}}  
{{- end}}
}

func New{{.StructName}}() *{{.StructName}} {
	return &{{.StructName}}{}
}
`

type CommandValues struct {
}

var commandTpl *template.Template

// GeneratorCommand 处理 JsonSchema 并生成代码
func (g *Generator) GeneratorCommand(schema *jsonschema.Schema, pkgName string, opts ...*Options) ([]byte, error) {
	opt := NewOptions(opts...)
	if opt.Template == nil {
		opt.Template = GetCommandTpl()
	}
	if opt.PackageName == "" {
		opt.PackageName = "command"
	}
	return g.Generate(opt.Template, schema, pkgName, opt.Values)
}

func GetCommandTpl() *template.Template {
	if commandTpl == nil {
		commandTpl = template.Must(template.New("model").Parse(CommandTemplate))
	}
	return commandTpl
}
