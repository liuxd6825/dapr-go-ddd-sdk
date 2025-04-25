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

type SchemaTemplate struct {
	tpl      *pongo2.Template
	tplSet   *pongo2.TemplateSet
	serverFs afero.Fs
	webFs    afero.Fs
}

func NewSchemaTemplate(tplSet *pongo2.TemplateSet, serverFs afero.Fs, webFs afero.Fs) *SchemaTemplate {
	return &SchemaTemplate{
		tplSet:   tplSet,
		serverFs: serverFs,
		webFs:    webFs,
	}
}

func (t *SchemaTemplate) Execute(ctx context.Context, sch *jsonschema.Schema, retValType string, valueVar string) (string, error) {
	data := map[string]any{
		"schema":     sch,
		"properties": GetAllProperties(sch),
		"retValType": retValType,
		"valueVar":   valueVar,
	}

	fsData, err := afero.ReadFile(t.webFs, "/@tpl/schema.tpl.html")
	if err != nil {
		return "", err
	}
	tpl, err := t.tplSet.FromBytes(fsData)
	if err != nil {
		panic("NewSchemaTemplate() " + err.Error())
	}

	str, err := tpl.Execute(data)
	return str, err
}

func (t *SchemaTemplate) Render(retValType string, valueVar string, htmlFile string, schemaFile string) (template.HTML, error) {
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
	if retValType != "" {
		retValType = strings.ToLower(retValType)
	}
	if valueVar == "" {
		valueVar = strings.ReplaceAll(schemaFile, "/", "_")
		valueVar = strings.ReplaceAll(valueVar, ".", "_")
	}
	h, err := t.Execute(context.Background(), sch, retValType, valueVar)
	if err != nil {
		return "", err
	}
	return template.HTML(h), nil
}
