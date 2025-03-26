package element

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/spf13/afero"
	"os"
)

type FsPkg interface {
	NewFs(fsName string) FsPkg
	Exists(fileName string, opts ...*fsopts.Options) bool
	Create(name string, opts ...*fsopts.Options) afero.File
	Rename(oldName, newName string, opts ...*fsopts.Options) error
	ReadFile(filename string, opts ...*fsopts.Options) []byte
	WriteJson(filename string, data any, opts ...*fsopts.Options)
	WriteFile(filename string, data any, opts ...*fsopts.Options)
	RemoveFile(filename string, opts ...*fsopts.Options)
	RemoveAll(name string, opts ...*fsopts.Options)
	Mkdir(name string, perm os.FileMode, opts ...*fsopts.Options)
	ReadDir(path string, opts ...*fsopts.Options) []*types.FileInfo
	ReadAllDir(path string, opts ...*fsopts.Options) []*types.FileInfo
	SortFileInfos(files []*types.FileInfo)
}
