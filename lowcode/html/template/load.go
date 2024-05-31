package template

import (
	"encoding/base64"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/spf13/afero"
	"io/ioutil"
	"strings"
)

// loadSchema
//
//	@Description: 通过js文件加载Schema
//	@param fs afero.Fs
//	@param pwd string 当前目录
//	@param schemasJsFile schemas js文件名称
//	@param schemaName string schema名称
//	@return *schema.Schema schema
//	@return error
func loadSchema(fs afero.Fs, pwd string, schemasFile string, schemaName string) (*schema.Schema, error) {
	bytes, err := readFile(fs, pwd, schemasFile)
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
//	@param pwd string 当前目录
//	@param uiSchemaJsFile string js文件名称
//	@param uiName string ui名称
//	@return *schema.UiSchema uiSchema
//	@return error error
func loadUiSchema(fs afero.Fs, sm *schema.Schema, pwd string, uiSchemaFile string, uiName string) (*schema.UiSchema, error) {
	if sm == nil {
		return nil, errors.New("schema instance is nil")
	}
	bytes, err := readFile(fs, pwd, uiSchemaFile)
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

// readFile
//
//	@Description: 读取文件内容
//	@param fs afero.Fs
//	@param pwdPath string 当前路径
//	@param filename string 要读取文件名称
//	@return []byte 文件内容
//	@return error
func readFile(fs afero.Fs, pwd string, filename string) ([]byte, error) {
	fileName := fileutils.AbsPath(pwd, filename)
	file, err := fs.Open(fileName)
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
