package metadata

import "github.com/liuxd6825/dapr-go-ddd-sdk/types"

//
// Base
// @Description:  元数据基类
//
type Base struct {
	Id      string   `json:"id,omitempty" bson:"id"`
	Name    string   `json:"name,omitempty" bson:"name"`
	Title   Title    `json:"title,omitempty" bson:"title"`
	Help    Help     `json:"help,omitempty" bson:"help"`
	Extends []string `json:"extends,omitempty" bson:"extends"`
}

//
// Help
// @Description:  帮助信息
//
type Help struct {
	CN string `json:"cn,omitempty" bson:"cn"`
	EN string `json:"en,omitempty" bson:"en"`
	DE string `json:"de,omitempty" bson:"de"`
}

//
// Title
// @Description: 标题
//
type Title struct {
	CN string `json:"cn,omitempty" bson:"cn"`
	EN string `json:"en,omitempty" bson:"en"`
	DE string `json:"de,omitempty" bson:"de"`
}

//
// Library
// @Description: 库
//
type Library struct {
	Namespace string           `json:"namespace,omitempty" bson:"namespace"`
	Models    map[string]Model `json:"models,omitempty" bson:"models"`
}

type Columns = map[string]Column
type DataType = types.DataType
type Fields = map[string]Field
type Tables = map[string]Grid
type Forms = map[string]Form
type Queries = map[string]Query

//
// Model
// @Description: 实体对象
//
type Model struct {
	Base
	Fields  Fields  `json:"fields,omitempty" bson:"fields"`
	Tables  Tables  `json:"tables,omitempty" bson:"tables"`
	Forms   Forms   `json:"forms,omitempty" bson:"forms"`
	Queries Queries `json:"queries,omitempty" bson:"queries"`
}

//
// Field
// @Description: 字段
//
type Field struct {
	PK       bool      `json:"pk,omitempty" bson:"pk"`
	DataType DataType  `json:"dataType,omitempty" bson:"data_type"`
	Length   int64     `json:"length,omitempty" bson:"length"`
	Form     FormItem  `json:"form,omitempty" bson:"form"`
	Column   Column    `json:"column,omitempty" bson:"column"`
	Query    QueryItem `json:"query,omitempty" bson:"query"`
}

//
// Foreign
// @Description: 字段外键
//
type Foreign struct {
	Base
	Model      string `json:"model,omitempty" bson:"model"`
	IdField    string `json:"idField,omitempty" bson:"id_field"`
	TitleField string `json:"titleField,omitempty" bson:"title_field"`
}

//
// Grid
// @Description: 表格
//
type Grid struct {
	Base
	Columns Columns `json:"columns,omitempty" bson:"columns"`
}

//
// Column
// @Description: 表格列
//
type Column struct {
	Field         string   `json:"field,omitempty,omitempty" bson:"field"`
	Width         int64    `json:"width,omitempty" bson:"width"`
	Editor        string   `json:"editor,omitempty" bson:"editor"`
	EditorOptions any      `json:"editorOptions,omitempty" bson:"editor_options"`
	Template      any      `json:"template,omitempty" bson:"template"`
	Align         string   `json:"align,omitempty" bson:"align"`
	CellsFormat   string   `json:"cellsFormat,omitempty" bson:"cells_format"`
	CellsAlign    string   `json:"cellsAlign,omitempty" bson:"cells_align"`
	ColumnGroup   string   `json:"columnGroup,omitempty" bson:"column_group"`
	Filter        string   `json:"filter,omitempty" bson:"filter"`
	Summary       []string `json:"summary,omitempty" bson:"summary"`
	Freeze        string   `json:"freeze,omitempty" bson:"freeze"`
	Visible       bool     `json:"visible,omitempty" bson:"visible"`

	AllowExport     bool `json:"allowExport,omitempty" bson:"allow_export"`
	AllowGroup      bool `json:"allowGroup,omitempty" bson:"allow_group"`
	AllowHide       bool `json:"allowHide,omitempty" bson:"allow_hide"`
	AllowSelect     bool `json:"allowSelect,omitempty" bson:"allow_select"`
	AllowEdit       bool `json:"allowEdit,omitempty" bson:"allow_edit"`
	AllowSort       bool `json:"allowSort,omitempty" bson:"allow_sort"`
	AllowHeaderEdit bool `json:"allowHeaderEdit,omitempty" bson:"allow_header_edit"`
	AllowFilter     bool `json:"allowFilter,omitempty" bson:"allow_filter"`
	AllowReorder    bool `json:"allowReorder,omitempty" bson:"allow_reorder"`
	AllowResize     bool `json:"allowResize,omitempty" bson:"allow_resize"`
	AllowNull       bool `json:"allowNull,omitempty" bson:"allow_null"`
}

//
// Editor
// @Description: 表格编辑器
//
type Editor struct {
	Title `json:"title,omitempty" bson:"title"`
	Help  `json:"help,omitempty" bson:"help"`
}

type QueryItems = map[string]QueryItem

//
// Query
// @Description: 查询
//
type Query struct {
	Base
}

//
// QueryItem
// @Description: 查询字段
//
type QueryItem struct {
	Base
	Field    string `json:"field,omitempty" bson:"field"`
	Operator string `json:"operator,omitempty" bson:"operator"`
	Value    any    `json:"value,omitempty" bson:"value"`
}

//
// Form
// @Description: 表单
//
type Form struct {
	Base
	Name   string              `json:"name,omitempty" bson:"name"`
	Cols   int64               `json:"cols,omitempty" bson:"cols"`
	Rows   int64               `json:"rows,omitempty" bson:"rows"`
	Fields map[string]FormItem `json:"fields,omitempty" bson:"fields"`
}

//
// FormItem
// @Description: 表单项
//
type FormItem struct {
	Field         string `json:"field,omitempty" bson:"field"`
	IsEdit        bool   `json:"isEdit,omitempty" bson:"is_edit"`
	Editor        string `json:"editor,omitempty" bson:"editor"`
	EditorOptions string `json:"editorOptions,omitempty" bson:"editor_options"`
	Cols          int64  `json:"cols,omitempty" bson:"cols"`
	Rows          int64  `json:"rows,omitempty" bson:"rows"`
	MaxValue      int64  `json:"maxValue,omitempty" bson:"max_value"`
	MinValue      int64  `json:"minValue,omitempty" bson:"min_value"`
	Script        string `json:"script,omitempty" bson:"script"`
}
