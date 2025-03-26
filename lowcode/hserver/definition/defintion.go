package definition

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
)

type Definition struct {
	fs        afero.Fs
	rootPath  string
	paramMaps map[string]*jsonschema.Schema
}

func NewDefinition(srcFs afero.Fs, rootPath string) (*Definition, error) {
	params := map[string]*jsonschema.Schema{}
	d := &Definition{fs: srcFs, paramMaps: params, rootPath: rootPath}
	if err := d.addParam(params, srcFs, rootPath); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Definition) addParam(params map[string]*jsonschema.Schema, srcFs afero.Fs, path string) error {
	fileInfos, err := afero.ReadDir(srcFs, path)
	if err != nil {
		return err
	}
	for _, fileInfo := range fileInfos {
		fileName := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			if err := d.addParam(params, srcFs, fileName); err != nil {
				return err
			}
		} else {
			paramsSch, err := d.newParam(srcFs, fileName)
			if err != nil {
				return err
			}
			fileName = strings.Replace(fileName, "\\", "/", -1) // 解决在window上目录格式问题
			//name := strings.Replace(fileName, d.rootPath+"/", "", 1)
			params[fileName] = paramsSch
		}
	}
	return nil
}

func (d *Definition) newParam(srcFs afero.Fs, fileName string) (*jsonschema.Schema, error) {
	data, err := fs.ReadFile(srcFs, fileName)
	if err != nil {
		return nil, err
	}
	if data != nil && len(data) > 0 {
		sch := schema.NewJsonSchemaWithJson(fileName, string(data))
		return sch, nil
	}
	return nil, nil
}

func (d *Definition) GetParam(key string) *jsonschema.Schema {
	params, ok := d.paramMaps[key]
	if !ok {
		return nil
	}
	return params
}
