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

func (m *MetaExtension) InitDBField(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "dbField")
	err := m.DBField.init(ctx, values)
	if m.DBField.Name == "id" {
		m.DBField.PrimaryKey = true
	}
	return err
}

func (m *MetaExtension) InitDBTable(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "dbTable")
	return m.DBTable.init(ctx, values)
}

func (m *MetaExtension) InitForm(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "form")
	return m.Form.init(ctx, values)
}

func (m *MetaExtension) InitColumn(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "column")
	return m.Column.init(ctx, values)
}

func (m *MetaExtension) InitLang(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "lang")
	return m.Lang.init(ctx, values)
}
func (m *MetaExtension) InitParam(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "param")
	return m.Param.init(ctx, values)
}

func (m *MetaExtension) InitDDD(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "ddd")
	return m.DDD.init(ctx, values)
}

func (m *MetaExtension) InitQuery(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "query")
	return m.Query.init(ctx, values)
}

func (m *MetaExtension) InitConvert(ctx *jsonschema.CompilerContext, meta map[string]any) error {
	values := getMapItem(meta, "convert")
	return m.Convert.init(ctx, values)
}

func (m *MetaExtension) Validate(ctx *jsonschema.ValidatorContext, v any) {

}
