package extensions

import (
	_ "embed"
	"github.com/liuxd6825/jsonschema/v6"
)

type MetaExtension struct {
	DBField *DBField   `json:"dbField,omitempty"`
	DBTable *DBTable   `json:"dbTable,omitempty"`
	Column  *Column    `json:"column,omitempty"`
	Query   *Query     `json:"query,omitempty"`
	Lang    *Lang      `json:"lang,omitempty"`
	Param   *HttpParam `json:"param,omitempty"`
}

//go:embed schema.json
var content string

const META_TAG_NAME = "meta"

func NewMetaExtension() *MetaExtension {
	return &MetaExtension{}
}

func (m *MetaExtension) InitDBField(vals map[string]any) error {
	m.DBField = &DBField{}
	return m.DBField.init(vals)
}

func (m *MetaExtension) InitDBTable(vals map[string]any) error {
	m.DBTable = &DBTable{}
	return m.DBTable.init(vals)
}

func (m *MetaExtension) InitColumn(vals map[string]any) error {
	m.Column = &Column{}
	return m.Column.init(vals)
}

func (m *MetaExtension) InitLang(vals map[string]any) error {
	m.Lang = &Lang{}
	return m.Lang.init(vals)
}

func (m *MetaExtension) InitQuery(vals map[string]any) error {
	m.Query = &Query{}
	return m.Query.init(vals)
}

func (m *MetaExtension) Validate(ctx *jsonschema.ValidatorContext, v any) {

}
