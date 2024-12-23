package pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"os"
)

type FileInfo struct {
	IsDir     bool        `json:"isDir"`
	Name      string      `json:"name"`
	Size      int64       `json:"size"`
	Path      string      `json:"path"`
	SubFiles  []*FileInfo `json:"subFiles"`
	SizeTitle string      `json:"sizeTitle"`
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
