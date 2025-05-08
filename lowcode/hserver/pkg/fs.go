package pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"os"
)

type FileInfo struct {
	IsDir     bool        `json:"isDir,omitempty"`
	Name      string      `json:"name,omitempty"`
	Size      int64       `json:"size,omitempty"`
	Path      string      `json:"path,omitempty"`
	SubFiles  []*FileInfo `json:"subFiles,omitempty"`
	SizeTitle string      `json:"sizeTitle,omitempty"`
}

type FsPkg interface {
	ReadFile(filename string, opts ...*fsopts.Options) []byte
	WriteFile(filename string, data any, opts ...*fsopts.Options)
	RemoveFile(filename string, opts ...*fsopts.Options)
	RemoveAll(name string, opts ...*fsopts.Options)
	Mkdir(name string, perm os.FileMode, opts ...*fsopts.Options)
	ReadDir(path string, opts ...*fsopts.Options) []*FileInfo
	ReadAllDir(path string, opts ...*fsopts.Options) []*FileInfo
}
