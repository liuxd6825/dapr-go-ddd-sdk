package builder

import (
	_ "embed"
	"github.com/flosch/pongo2/v6"
)

type Template struct {
	tplData  []byte
	template *pongo2.Template
}

func NewTemplate(tplData []byte) *Template {
	tpl := newTemplate(tplData)
	return &Template{
		tplData:  tplData,
		template: tpl,
	}
}

func newTemplate(tplData []byte) *pongo2.Template {
	tpl, err := pongo2.FromBytes(tplData)
	if err != nil {
		panic(err)
	}
	return tpl
}

func (t *Template) Execute(data pongo2.Context) (string, error) {
	str, err := t.template.Execute(data)
	if err != nil {
		return "", err
	}
	return str, nil
}
