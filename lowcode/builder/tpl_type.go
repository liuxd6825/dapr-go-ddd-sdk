package builder

import (
	_ "embed"
	"strings"
)

//go:embed tpl/form.tpl.html
var formFile []byte
var formTpl *Tpl

//go:embed tpl/grid.tpl.html
var gridFile []byte
var gridTpl *Tpl

type TplType int

const (
	TplTypeNull TplType = iota
	TplTypeForm
	TplTypeGrid
)

func init() {
	formTpl = NewTpl(formFile)
	gridTpl = NewTpl(gridFile)
}

func (t TplType) String() string {
	switch t {
	case TplTypeForm:
		return "form"
	case TplTypeGrid:
		return "grid"
	default:
		return ""
	}
}

func GetTplType(name string) TplType {
	n := strings.ToLower(name)
	switch n {
	case "form":
		return TplTypeForm
	case "grid":
		return TplTypeGrid
	default:
		return TplTypeNull
	}

}
func (t TplType) GetTpl() *Tpl {
	switch t {
	case TplTypeForm:
		return formTpl
	case TplTypeGrid:
		return gridTpl
	default:
		return nil
	}
}
