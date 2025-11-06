package builder

import (
	_ "embed"
	"strings"
)

//go:embed tpl/form.tpl.html
var formFile []byte
var formTpl *Template

//go:embed tpl/grid.tpl.html
var gridFile []byte
var gridTpl *Template

//go:embed tpl/service.tpl.html
var serviceFile []byte
var serviceTpl *Template

type TemplateType string

const (
	TplTypeNull    TemplateType = "null"
	TplTypeForm    TemplateType = "form"
	TplTypeGrid    TemplateType = "grid"
	TplTypeService TemplateType = "service"
)

func init() {
	formTpl = NewTemplate(formFile)
	gridTpl = NewTemplate(gridFile)
	serviceTpl = NewTemplate(serviceFile)
}

func NewTemplateType(name string) TemplateType {
	n := strings.ToLower(name)
	switch n {
	case "form":
		return TplTypeForm
	case "grid":
		return TplTypeGrid
	case "service":
		return TplTypeService
	default:
		return TplTypeNull
	}
}

func (t TemplateType) String() string {
	return string(t)
}

func (t TemplateType) GetTemplate() *Template {
	switch t {
	case TplTypeForm:
		return formTpl
	case TplTypeGrid:
		return gridTpl
	case TplTypeService:
		return serviceTpl
	default:
		return nil
	}
}
