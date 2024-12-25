package builder

import _ "embed"

//go:embed tpl/form.tpl.html
var formFile []byte
var formTpl *Tpl

//go:embed tpl/grid.tpl.html
var gridFile []byte
var gridTpl *Tpl

type TplType int

const (
	TplTypeForm = TplType(iota)
	TplTypeGrid
)

func init() {
	formTpl = NewTpl(formFile)
	gridTpl = NewTpl(gridFile)
}

func (t TplType) String() string {
	switch t {
	case TplTypeForm:
		return "Form"
	case TplTypeGrid:
		return "Grid"
	}
	return ""
}

func (t TplType) GetTpl() *Tpl {
	switch t {
	case TplTypeForm:
		return formTpl
	case TplTypeGrid:
		return gridTpl
	}
	return nil
}
