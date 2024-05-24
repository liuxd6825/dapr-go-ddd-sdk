package template

import (
	"encoding/base64"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/spf13/afero"
	"io/ioutil"
	"strings"
)

// loadSchema
//
//	@Description: 通过js文件加载Schema
//	@param fs afero.Fs
//	@param schemasJsFile schemas js文件名称
//	@param schemaName string schema名称
//	@return *schema.Schema schema
//	@return error
func loadSchema(fs afero.Fs, schemasFile string, schemaName string) (*schema.Schema, error) {
	bytes, err := readFile(fs, schemasFile)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(schemasFile, ".js") {
		schemas := schema.NewSchemasLibrary()
		if err = schemas.LoadJavaScript(bytes); err != nil {
			return nil, err
		}
		s, err := schemas.Get(schemaName)
		return s, err
	}

	s, err := schema.NewSchemaWithJson(string(bytes))
	return s, err
}

// loadUiSchema
//
//	@Description: 通过js文件加载UiSchema
//	@param fs afero.Fs
//	@param sm *schema.Schema
//	@param uiSchemaJsFile string js文件名称
//	@param uiName string ui名称
//	@return *schema.UiSchema uiSchema
//	@return error error
func loadUiSchema(fs afero.Fs, sm *schema.Schema, uiSchemaFile string, uiName string) (*schema.UiSchema, error) {
	if sm == nil {
		return nil, errors.New("schema instance is nil")
	}
	bytes, err := readFile(fs, uiSchemaFile)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(uiSchemaFile, ".js") {
		uiSchemas := schema.NewUiSchemasLibrary()
		if err = uiSchemas.LoadJavaScript(bytes); err != nil {
			return nil, err
		}
		ui, err := uiSchemas.Get(uiName)
		if ui != nil {
			if err = ui.Init(sm); err != nil {
				return nil, err
			}
		}
		return ui, err
	}
	ui, err := schema.NewUiSchemaWithJson(sm, string(bytes))
	return ui, err
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
