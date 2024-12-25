package builder

import (
	_ "embed"
	"github.com/flosch/pongo2/v6"
)

type Tpl struct {
	tplData []byte
	tpl     *pongo2.Template
}

func NewTpl(tplData []byte) *Tpl {
	tpl := newTpl(tplData)
	return &Tpl{
		tplData: tplData,
		tpl:     tpl,
	}
}

func newTpl(tplData []byte) *pongo2.Template {
	tpl, err := pongo2.FromBytes(tplData)
	if err != nil {
		panic(err)
	}
	return tpl
}

func (t *Tpl) Execute(data pongo2.Context) (string, error) {
	str, err := t.tpl.Execute(data)
	if err != nil {
		return "", err
	}
	return str, nil
}
