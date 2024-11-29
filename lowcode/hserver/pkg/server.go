package pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/spf13/afero"
)

type Server interface {
	GetFsm() *fs_pkg.FsManager
	GetSrcFs() afero.Fs
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
}
