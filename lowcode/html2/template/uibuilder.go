package template

import (
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"io"
	"strings"
)

type UiBuilder struct {
	schema   *schema.Schema
	uiSchema *schema.UiSchema
}

type BuildResult struct {
	sb strings.Builder
}

type BuildGrid struct {
	schema.UiLayout
	Items []BuildRow `json:"rows"`
}

type BuildRow struct {
	schema.UiRow
	Items []BuildCol `json:"cols"`
}

type BuildCol struct {
	schema.Property
	UiProperty *schema.UiProperty
}

func NewUiBuilder(schema *schema.Schema, uiSchema *schema.UiSchema) *UiBuilder {
	return &UiBuilder{
		schema:   schema,
		uiSchema: uiSchema,
	}
}

func (b *UiBuilder) init(schema *schema.Schema, uiSchema *schema.UiSchema) error {
	b.schema = schema
	b.uiSchema = uiSchema
	return nil
}

func (b *UiBuilder) LoadBytes(schemaByte []byte, uiSchemaByte []byte) error {
	s, err := schema.NewSchema("", strings.NewReader(string(schemaByte)))
	if err != nil {
		return err
	}

	uiSchema, err := schema.NewUiSchema(s, strings.NewReader(string(uiSchemaByte)))
	if err != nil {
		return err
	}

	return b.init(s, uiSchema)
}

func (b *UiBuilder) LoadFile(schemaFile string, uiSchemaFile string) error {
	s, err := schema.NewSchemaFile(schemaFile)
	if err != nil {
		return err
	}

	uiSchema, err := schema.NewUiSchemaFile(s, uiSchemaFile)
	if err != nil {
		return err
	}

	return b.init(s, uiSchema)
}

func (b *UiBuilder) Build(tplFile string, writer io.Writer) error {
	tpl, err := pongo2.FromFile(tplFile)
	if err != nil {
		return err
	}
	return b.build(tpl, writer)
}

func (b *UiBuilder) BuildBytes(content []byte, writer io.Writer) error {
	tpl, err := pongo2.FromBytes(content)
	if err != nil {
		return err
	}

	return b.build(tpl, writer)
}

func (b *UiBuilder) BuildString(content string, writer io.Writer) error {
	tpl, err := pongo2.FromString(content)
	if err != nil {
		return err
	}
	return b.build(tpl, writer)
}

func (b *UiBuilder) build(tpl *pongo2.Template, writer io.Writer) error {
	errs := errors.NewErrors()
	if b.schema == nil {
		errs.AddFormat("parameter schema is nil")
	}

	if b.uiSchema == nil {
		errs.AddFormat("parameter uiSchema is nil")
	} else if b.uiSchema.Layout == nil {
		errs.AddFormat("parameter uiSchema.layout is nil")
	} else if b.uiSchema.Layout.Rows == nil {
		errs.AddFormat("parameter uiSchema.layout.Rows  is nil")
	}

	if errs.HasError() {
		return errs.NewError()
	}

	ctx := pongo2.Context{"schema": b.schema, "uischema": b.uiSchema, "grid": b.uiSchema.Layout}
	return tpl.ExecuteWriter(ctx, writer)
}
