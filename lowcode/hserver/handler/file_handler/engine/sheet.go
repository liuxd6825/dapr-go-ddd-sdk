package engine

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/spf13/afero"
	"html/template"
	"strings"
)

//go:embed "tpl/sheet.tpl.html"
var _sheetBytes []byte // 嵌入为字节切片

type SheetTemplate struct {
	tpl      *pongo2.Template
	serverFs afero.Fs
	webFs    afero.Fs
}

type Property struct {
	*jsonschema.Schema
	meta *schema.MetaExtension
}

func NewSheetTemplate(serverFs afero.Fs, webFs afero.Fs) *SheetTemplate {
	tpl, err := pongo2.FromBytes(_sheetBytes)
	if err != nil {
		panic("NewSheetTemplate() " + err.Error())
	}
	return &SheetTemplate{
		tpl:      tpl,
		serverFs: serverFs,
		webFs:    webFs,
	}
}

func (t *SheetTemplate) Execute(ctx context.Context, sch *jsonschema.Schema) (string, error) {
	properties := make([]*Property, 0)
	for _, item := range sch.GetSortProperties() {
		prop := Property{
			Schema: item,
		}
		properties = append(properties, &prop)
	}
	data := map[string]any{
		"schema":     sch,
		"properties": properties,
	}

	str, err := t.tpl.Execute(data)
	return str, err
}

type SheetOptions struct {
	Schema     *jsonschema.Schema
	SchemaFile string
}

func (p *Property) IsIdField() bool {
	if strings.HasSuffix(strings.ToLower(p.Name()), "id") {
		return true
	}
	return false
}

func (p *Property) Meta() *schema.MetaExtension {
	if p.meta == nil {
		p.meta = schema.GetMetaExtension(p.Schema)
	}
	if p.meta == nil {
		p.meta = schema.NewMetaExtension()
	}
	return p.meta
}

func (p *Property) Readonly() bool {
	return p.Meta().DBField.Readonly()
}

func (p *Property) Type() string {
	return "text"
	if p.Types.Contains(jsonschema.JsonType_BooleanType) {
		return "boolean"
	} else if p.Types.Contains(jsonschema.JsonType_NumberType) {
		return "number"
	} else if p.Types.Contains(jsonschema.JsonType_StringType) {
		return "string"
	} else if p.Types.Contains(jsonschema.JsonType_IntegerType) {
		return "integer"
	} else if p.Types.Contains(jsonschema.JsonType_DateType) {
		return "date"
	} else if p.Types.Contains(jsonschema.JsonType_DateTimeType) {
		return "datetime"
	}
	return ""
}

func (t *SheetTemplate) Render(htmlFile string, schemaFile string) (template.HTML, error) {
	/*
		htmlFile, _ := maputils.GetString(opts, "html", "")
		schemaFile, _ := maputils.GetString(opts, "schema", "")
	*/
	if htmlFile == "" && schemaFile == "" {
		return template.HTML(""), nil
	}
	if htmlFile != "" {
		if ok, _ := afero.Exists(t.webFs, htmlFile); ok {
			h, err := afero.ReadFile(t.webFs, htmlFile)
			if err != nil {
				return "", err
			}
			return template.HTML(h), nil
		}
	}

	var sch *jsonschema.Schema
	if schemaFile != "" {
		data, err := afero.ReadFile(t.serverFs, schemaFile)
		if err != nil {
			return "", err
		}
		sch = schema.NewJsonSchemaWithBytes(schemaFile, data)
	}
	if sch == nil {
		return "", fmt.Errorf("no schema found in %s", schemaFile)
	}
	h, err := t.Execute(context.Background(), sch)
	if err != nil {
		return "", err
	}
	return template.HTML(h), nil
}
