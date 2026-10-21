package fs_pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
)

func NewFsWriteModel() *fspkg.FsWriteModel {
	return fspkg.NewFsWriteModel()
}

// NewFsPkg
//
//	@Description: 创建文件系统
//	@param cfg
//	@return *FsPkg
//	@return error
func NewFsPkg(cfg *envs.Env, fsName string) (fspkg.IFsPkg, error) {
	return fspkg.NewFsPkg(cfg, fsName)
}
