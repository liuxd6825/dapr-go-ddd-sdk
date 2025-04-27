package schema

import (
	_ "embed"
	"github.com/liuxd6825/jsonschema/v6"
)

type MetaExtension struct {
	DBField *DBField   `json:"dbField,omitempty"`
	DBTable *DBTable   `json:"dbTable,omitempty"`
	Form    *Form      `json:"form,omitempty"`
	Column  *Column    `json:"column,omitempty"`
	Query   *Query     `json:"query,omitempty"`
	Lang    *Lang      `json:"lang,omitempty"`
	Param   *HttpParam `json:"param,omitempty"`
	Convert *Convert   `json:"convert,omitempty"` // 数据转换器的名称
	DDD     *DDD       `json:"ddd,omitempty"`
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

func (m *MetaExtension) InitDBField(vals map[string]any) error {
	m.DBField = &DBField{}
	return m.DBField.init(vals)
}

func (m *MetaExtension) InitDBTable(vals map[string]any) error {
	m.DBTable = &DBTable{}
	return m.DBTable.init(vals)
}

func (m *MetaExtension) InitForm(vals map[string]any) error {
	m.Form = &Form{}
	return m.Form.init(vals)
}

func (m *MetaExtension) InitColumn(vals map[string]any) error {
	m.Column = &Column{}
	return m.Column.init(vals)
}

func (m *MetaExtension) InitLang(vals map[string]any) error {
	m.Lang = &Lang{}
	return m.Lang.init(vals)
}
func (m *MetaExtension) InitParam(vals map[string]any) error {
	m.Param = &HttpParam{}
	return m.Param.init(vals)
}

func (m *MetaExtension) InitDDD(vals map[string]any) error {
	m.DDD = &DDD{}
	return m.DDD.init(vals)
}

func (m *MetaExtension) InitQuery(vals map[string]any) error {
	m.Query = &Query{}
	return m.Query.init(vals)
}

func (m *MetaExtension) InitConvert(vals map[string]any) error {
	m.Convert = &Convert{}
	return m.Convert.init(vals)
}

func (m *MetaExtension) Validate(ctx *jsonschema.ValidatorContext, v any) {

}
