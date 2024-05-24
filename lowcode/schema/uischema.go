package schema

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"io"
	"os"
	"strings"
)

type UiSchema struct {
	Field      string       `json:"field,omitempty"`
	Layout     *UiLayout    `json:"layout,omitempty"`
	Title      string       `json:"title,omitempty"`
	Properties UiProperties `json:"properties,omitempty"`
	schema     *Schema
}

type UiLayout struct {
	TitleWidth string `json:"titleWidth,omitempty"`
	Rows       UiRows `json:"rows,omitempty"`
}

type UiRows []*UiRow

type UiRow struct {
	Row    int
	Height string   `json:"height,omitempty"`
	Cols   []*UiCol `json:"cols,omitempty"`
}

type UiCol struct {
	Name       string      `json:"name,omitempty"`
	Span       int         `json:"span,omitempty"`
	Required   bool        `json:"required,omitempty"`
	UiProperty *UiProperty `json:"-"`
	Property   *Property   `json:"-"`
	schema     *Schema     `json:"-"`
}

type UiProperties map[string]*UiProperty

type UiProperty struct {
	Widget                      string             `json:"widget,omitempty"`
	Options                     *UiPropertyOptions `json:"options,omitempty"`
	Readonly                    string             `json:"readonly,omitempty"`
	Autofocus                   string             `json:"autofocus,omitempty"`
	Autocomplete                string             `json:"autocomplete,omitempty"`
	Description                 string             `json:"description,omitempty"`
	EmptyValue                  string             `json:"emptyValue,omitempty"`
	Placeholder                 string             `json:"placeholder,omitempty"`
	EnableMarkdownInDescription string             `json:"enableMarkdownInDescription,omitempty"`
}

type UiPropertyOptions struct {
	InputType string `json:"inputType,omitempty"`
}

func NewUiSchemaFile(schema *Schema, pathFile string) (*UiSchema, error) {
	// 打开 JSON 文件
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error opening file: %v", err))
	}
	defer file.Close()

	return NewUiSchema(schema, file)
}

func NewUiSchemaWithJson(schema *Schema, jsonStr string) (*UiSchema, error) {
	var reader = strings.NewReader(jsonStr)
	return NewUiSchema(schema, reader)
}

func NewUiSchema(schema *Schema, reader io.Reader) (*UiSchema, error) {
	var ui UiSchema
	var err error
	// 解码 JSON 数据到结构体中
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&ui); err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding JSON: %v", err))
	}
	if err = ui.Init(schema); err != nil {
		return nil, err
	}
	return &ui, err
}

func (s *UiSchema) Init(schema *Schema) error {
	if schema == nil {
		return errors.New("schema is nil")
	}
	s.schema = schema
	if s.Layout == nil {
		return errors.New("ui:layout is required")
	}
	errs := errors.NewErrors()
	for iRow, row := range s.Layout.Rows {
		for iCol, col := range row.Cols {
			name := col.Name
			if name == "" {
				errs.AddFormat("row[%d].col[%d].name is required", iRow, iCol)
				continue
			}

			col.schema = schema
			if col.Span < 0 {
				col.Span = 1
			}

			// 设置col.Required是否必填
			col.Required = false
			for _, key := range schema.Required {
				if col.Name == key {
					col.Required = true
					break
				}
			}

			// 设置col.Property和col.UiProperty
			if prop, ok := schema.Properties[name]; ok {
				col.Property = prop
			} else {
				errs.AddFormat("schema.properties.%s not found", name)
			}
			if prop, ok := s.Properties[name]; ok {
				col.UiProperty = prop
			} else {
				errs.AddFormat("uiSchema.properties.%s not found", name)
			}
		}
	}
	if errs.HasError() {
		return errs
	}
	return nil
}

func (g *UiLayout) GetTitleWidth() string {
	if g.TitleWidth == "" {
		return "200px"
	}
	return g.TitleWidth
}

func (c *UiRow) GetStyle() string {
	sb := strings.Builder{}
	if c.Height != "" {
		sb.WriteString(fmt.Sprintf("min-height:%s;", c.Height))
	}
	return sb.String()
}

func (c *UiCol) GetStyle(width string) string {
	sb := strings.Builder{}
	if c.Span > 0 {
		sb.WriteString(fmt.Sprintf("flex:0 0 calc(%v\\%% - %s", c.Span*100/24, width))
	}
	return sb.String()
}

func (s *UiSchema) ToJson() string {
	bs, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}
