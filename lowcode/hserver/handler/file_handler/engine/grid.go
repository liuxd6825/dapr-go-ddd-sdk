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
)

type AgGridTemplate struct {
	tpl      *pongo2.Template
	tplSet   *pongo2.TemplateSet
	serverFs afero.Fs
	webFs    afero.Fs
}

func NewAgGridTemplate(tplSet *pongo2.TemplateSet, serverFs afero.Fs, webFs afero.Fs) *AgGridTemplate {
	return &AgGridTemplate{
		tplSet:   tplSet,
		serverFs: serverFs,
		webFs:    webFs,
	}
}

func (t *AgGridTemplate) Execute(ctx context.Context, sch *jsonschema.Schema) (string, error) {
	data := map[string]any{
		"schema":     sch,
		"properties": GetAllProperties(sch),
	}

	fsData, err := afero.ReadFile(t.webFs, "/@tpl/grid.tpl.html")
	if err != nil {
		return "", err
	}
	tpl, err := t.tplSet.FromBytes(fsData)
	if err != nil {
		panic("NewAgGridTemplate() " + err.Error())
	}

	str, err := tpl.Execute(data)
	return str, err
}

func (t *AgGridTemplate) Render(htmlFile string, schemaFile string) (template.HTML, error) {
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
