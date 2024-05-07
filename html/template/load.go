package template

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
)

func loadSchema(loader loader.Loader, schemaFile string) (*schema.Schema, error) {
	bytes, err := loader.GetFile(schemaFile)
	if err != nil {
		return nil, err
	}
	return schema.NewSchemaString(string(bytes))
}

func loadUiSchema(loader loader.Loader, sm *schema.Schema, url string) (*schema.UiSchema, error) {
	bytes, err := loader.GetFile(url)
	if err != nil {
		return nil, err
	}
	return schema.NewUiSchemaString(sm, string(bytes))
}
