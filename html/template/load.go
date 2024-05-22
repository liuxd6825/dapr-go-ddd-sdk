package template

import (
	"encoding/base64"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
	"github.com/spf13/afero"
	"io/ioutil"
)

func loadSchema(fs afero.Fs, schemaFile string) (*schema.Schema, error) {
	bytes, err := readFile(fs, schemaFile)
	if err != nil {
		return nil, err
	}
	return schema.NewSchemaString(string(bytes))
}

func loadUiSchema(fs afero.Fs, sm *schema.Schema, uiSchemaFile string) (*schema.UiSchema, error) {
	bytes, err := readFile(fs, uiSchemaFile)
	if err != nil {
		return nil, err
	}
	return schema.NewUiSchemaString(sm, string(bytes))
}

func readFile(fs afero.Fs, filename string) ([]byte, error) {
	file, err := fs.Open(filename)
	if err != nil {
		return nil, err
	}
	context, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if fs.Name() == giteafs.Name() {
		context, err = base64.StdEncoding.DecodeString(string(context))
	}
	return context, err
}
