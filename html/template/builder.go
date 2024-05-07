package template

import (
	"errors"
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
	"io"
	"strings"
)

type Builder struct {
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

func NewBuilder(schema *schema.Schema, uiSchema *schema.UiSchema) *Builder {
	return &Builder{
		schema:   schema,
		uiSchema: uiSchema,
	}
}

func (b *Builder) init(schema *schema.Schema, uiSchema *schema.UiSchema) error {
	b.schema = schema
	b.uiSchema = uiSchema
	return nil
}

func (b *Builder) LoadBytes(schemaByte []byte, uiSchemaByte []byte) error {
	s, err := schema.NewSchema(strings.NewReader(string(schemaByte)))
	if err != nil {
		return err
	}

	uiSchema, err := schema.NewUiSchema(s, strings.NewReader(string(uiSchemaByte)))
	if err != nil {
		return err
	}

	return b.init(s, uiSchema)
}

func (b *Builder) LoadFile(schemaFile string, uiSchemaFile string) error {
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

func (b *Builder) Build(tplFile string, writer io.Writer) error {
	tpl, err := pongo2.FromFile(tplFile)
	if err != nil {
		return err
	}
	return b.build(tpl, writer)
}

func (b *Builder) BuildBytes(content []byte, writer io.Writer) error {
	tpl, err := pongo2.FromBytes(content)
	if err != nil {
		return err
	}

	return b.build(tpl, writer)
}

func (b *Builder) BuildString(content string, writer io.Writer) error {
	tpl, err := pongo2.FromString(content)
	if err != nil {
		return err
	}
	return b.build(tpl, writer)
}

func (b *Builder) build(tpl *pongo2.Template, writer io.Writer) error {
	if b.uiSchema == nil || b.uiSchema.Layout == nil || b.uiSchema.Layout.Rows == nil {
		return errors.New("")
	}
	ctx := pongo2.Context{"schema": b.schema, "uischema": b.uiSchema, "grid": b.uiSchema.Layout}
	return tpl.ExecuteWriter(ctx, writer)
}
