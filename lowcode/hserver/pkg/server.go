package pkg

import "github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"

type Server interface {
	GetFs() *fs_pkg.FsManager
	GetSrcPath() string
}
