package definition

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/xtype"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/spf13/afero"
	iofs "io/fs"
)

type Definition struct {
	fs       afero.Fs
	requests map[string]xtype.ParamsType `json:"requests"`
}

func NewDefinition(srcFs afero.Fs, path string) (*Definition, error) {
	requests := map[string]xtype.ParamsType{}
	err := afero.Walk(srcFs, path, func(path string, info iofs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 打印文件路径（忽略目录）
		if !info.IsDir() {
			data, err := fs.ReadFile(srcFs, path)
			if err != nil {
				return err
			}
			if data != nil && len(data) > 0 {
				var paramsType xtype.ParamsType
				if err = jsonutils.Unmarshal(data, &paramsType); err != nil {
					panic(fmt.Sprintf(" loading %s  error: %s", path, err.Error()))
				}
				requests[info.Name()] = paramsType
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Definition{fs: srcFs, requests: requests}, nil
}

func (d *Definition) GetParamsType(key string) xtype.ParamsType {
	params, ok := d.requests[key]
	if !ok {
		return nil
	}
	return params
}
