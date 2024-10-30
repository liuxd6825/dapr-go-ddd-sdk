package server

import (
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
)

type Template struct {
	cfg common.IEnvConfig
}

func NewTemplate(cfg common.IEnvConfig) *Template {
	return &Template{cfg: cfg}
}

func (e *Template) RenderFile(filename string, data map[string]any) string {
	bytes := GetFsManger().ReadFile(filename, "")
	txt := e.renderBytes(bytes, data)
	return *txt
}

func (e *Template) RenderString(str string, data map[string]any) string {
	tpl, err := pongo2.FromString(str)
	if err != nil {
		panic(err)
	}
	txt := e.render(tpl, err, data)
	return *txt
}

func (e *Template) renderBytes(bytes []byte, data map[string]any) *string {
	tpl, err := pongo2.FromBytes(bytes)
	if err != nil {
		panic(err)
	}
	return e.render(tpl, err, data)
}

func (e *Template) render(tpl *pongo2.Template, err error, data map[string]any) *string {
	if err != nil {
		panic(err)
	}
	context, err := tpl.Execute(data)
	if err != nil {
		panic(err)
	}
	return &context
}
