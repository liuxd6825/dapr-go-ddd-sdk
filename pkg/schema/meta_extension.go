package schema

import (
	_ "embed"
	"github.com/liuxd6825/jsonschema/v6"
)

type MetaExtension struct {
	DBField *DBField           `json:"dbField,omitempty"`
	DBTable *DBTable           `json:"dbTable,omitempty"`
	Form    *Form              `json:"form,omitempty"`
	Column  *Column            `json:"column,omitempty"`
	Query   *Query             `json:"query,omitempty"`
	Lang    *Lang              `json:"lang,omitempty"`
	Param   *HttpParam         `json:"param,omitempty"`
	Convert *Convert           `json:"convert,omitempty"` // 数据转换器的名称
	DDD     *DDD               `json:"ddd,omitempty"`
	sch     *jsonschema.Schema `json:"-"`
}

//go:embed schema.json
var content string

const META_TAG_NAME = "meta"

func NewMetaExtension() *MetaExtension {
	return &MetaExtension{
		DBField: NewDBField(),
		DBTable: NewDBTable(),
		Form:    NewForm(),
		Column:  NewColumn(),
		Query:   NewQuery(),
		Lang:    NewLang(),
		Param:   NewHttpParam(),
		Convert: NewConvert(),
		DDD:     NewDDD(),
	}
}

func (m *MetaExtension) TagName() string {
	return META_TAG_NAME
}

func (m *MetaExtension) InitDBField(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.DBField = &DBField{}
	return m.DBField.init(ctx, vals)
}

func (m *MetaExtension) InitDBTable(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.DBTable = &DBTable{}
	return m.DBTable.init(ctx, vals)
}

func (m *MetaExtension) InitForm(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.Form = &Form{}
	return m.Form.init(ctx, vals)
}

func (m *MetaExtension) InitColumn(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.Column = &Column{}
	return m.Column.init(ctx, vals)
}

func (m *MetaExtension) InitLang(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.Lang = &Lang{}
	return m.Lang.init(ctx, vals)
}
func (m *MetaExtension) InitParam(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.Param = &HttpParam{}
	return m.Param.init(ctx, vals)
}

func (m *MetaExtension) InitDDD(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.DDD = &DDD{}
	return m.DDD.init(ctx, vals)
}

func (m *MetaExtension) InitQuery(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.Query = &Query{}
	return m.Query.init(ctx, vals)
}

func (m *MetaExtension) InitConvert(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	m.Convert = &Convert{}
	return m.Convert.init(ctx, vals)
}

func (m *MetaExtension) Validate(ctx *jsonschema.ValidatorContext, v any) {

}
