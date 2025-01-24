package definition

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/spf13/afero"
	"path/filepath"
	"strings"
)

type Definition struct {
	fs       afero.Fs
	rootPath string
	params   map[string]common.ParamsType
}

func NewDefinition(srcFs afero.Fs, rootPath string) (*Definition, error) {
	params := map[string]common.ParamsType{}
	d := &Definition{fs: srcFs, params: params, rootPath: rootPath}
	if err := d.addParamsType(params, srcFs, rootPath); err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Definition) addParamsType(params map[string]common.ParamsType, srcFs afero.Fs, path string) error {
	fileInfos, err := afero.ReadDir(srcFs, path)
	if err != nil {
		return err
	}
	for _, fileInfo := range fileInfos {
		fileName := filepath.Join(path, fileInfo.Name())
		if fileInfo.IsDir() {
			if err := d.addParamsType(params, srcFs, fileName); err != nil {
				return err
			}
		} else {
			paramsType, err := d.newParamsType(srcFs, fileName)
			if err != nil {
				return err
			}
			name := strings.Replace(fileName, d.rootPath+"/", "", 1)
			params[name] = paramsType
		}
	}
	return nil
}

func (d *Definition) newParamsType(srcFs afero.Fs, fileName string) (common.ParamsType, error) {
	data, err := fs.ReadFile(srcFs, fileName)
	if err != nil {
		return nil, err
	}
	if data != nil && len(data) > 0 {
		var paramsType common.ParamsType
		if err = jsonutils.Unmarshal(data, &paramsType); err != nil {
			return nil, errors.New(" loading %s  error: %s", fileName, err.Error())
		}
		return paramsType, nil
	}
	return nil, nil
}

func (d *Definition) GetParamsType(key string) common.ParamsType {
	params, ok := d.params[key]
	if !ok {
		return nil
	}
	return params
}
