package test

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"os"
)

// NewFsManager
//
//	@Description:创建一个测试对象
//	@return *fs.Manager
//	@return error
func NewFsManager(rootPath string) (*fs.Manager, error) {
	// 获取当前工作目录
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fsm := fs.NewManager()
	fsCfg := localfs.Config{
		Name: "file",
		Path: path + rootPath,
	}
	lfs, err := localfs.NewFs(fsCfg)
	if err != nil {

		return nil, err
	}

	fsm.Add("file", lfs)
	fsm.DefaultFsName = "file"
	return fsm, nil
}
