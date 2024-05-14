package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type UiSchema struct {
	Field      string       `json:"ui:field,omitempty"`
	Layout     *UiLayout    `json:"ui:layout,omitempty"`
	Title      string       `json:"ui:title,omitempty"`
	Properties UiProperties `json:"ui:properties,omitempty"`
	schema     *Schema
}

type UiLayout struct {
	TitleWidth string `json:"ui:titleWidth,omitempty"`
	Rows       UiRows `json:"ui:rows,omitempty"`
}

type UiRows []*UiRow

type UiRow struct {
	Row    int
	Height string   `json:"ui:height,omitempty"`
	Cols   []*UiCol `json:"ui:cols,omitempty"`
}

type UiCol struct {
	Name       string      `json:"name,omitempty"`
	Span       int         `json:"span,omitempty"`
	UiProperty *UiProperty `json:"uiproperty,omitempty"`
	Property   *Property   `json:"property,omitempty"`
	schema     *Schema     `json:"validate,omitempty"`
	Required   bool        `json:"required,omitempty"`
}

type UiProperties map[string]*UiProperty

type UiProperty struct {
	Widget                      string             `json:"ui:widget,omitempty"`
	Options                     *UiPropertyOptions `json:"ui:options,omitempty"`
	Readonly                    string             `json:"ui:readonly,omitempty"`
	Autofocus                   string             `json:"ui:autofocus,omitempty"`
	Autocomplete                string             `json:"ui:autocomplete,omitempty"`
	Description                 string             `json:"ui:description,omitempty"`
	EmptyValue                  string             `json:"ui:emptyValue,omitempty"`
	Placeholder                 string             `json:"ui:placeholder,omitempty"`
	EnableMarkdownInDescription string             `json:"ui:enableMarkdownInDescription,omitempty"`
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

func NewUiSchemaString(schema *Schema, jsonStr string) (*UiSchema, error) {
	var reader = strings.NewReader(jsonStr)
	return NewUiSchema(schema, reader)
}

func NewUiSchema(schema *Schema, reader io.Reader) (*UiSchema, error) {
	var data UiSchema
	var err error
	// 解码 JSON 数据到结构体中
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&data); err != nil {
		return nil, errors.New(fmt.Sprintf("Error decoding JSON: %v", err))
	}
	if err = data.Init(schema); err != nil {
		return nil, err
	}
	return &data, err
}

func (s *UiSchema) Init(schema *Schema) error {
	s.schema = schema
	for _, row := range s.Layout.Rows {
		for _, col := range row.Cols {
			col.schema = schema
			if col.Span < 0 {
				col.Span = 1
			}
			col.Required = false
			for _, key := range schema.Required {
				if col.Name == key {
					col.Required = true
					break
				}
			}
			name := col.Name
			if prop, ok := schema.Properties[name]; ok {
				col.Property = prop
			}
			if prop, ok := s.Properties[name]; ok {
				col.UiProperty = prop
			}
		}
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
		sb.WriteString(fmt.Sprintf("flex:0 0 calc(%v\\% - %s", c.Span*100/24, width))
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
