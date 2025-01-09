package hserver

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"github.com/spf13/afero"
)

// SrcFs
// @Description: source code file system config
type SrcFs struct {
	FileFs afero.Fs
	HttpFs afero.Fs
}

func NewSrcFs(fileFs, httpFs afero.Fs) *SrcFs {
	return &SrcFs{FileFs: fileFs, HttpFs: httpFs}
}

func (c *SrcFs) ToMap() map[string]afero.Fs {
	return map[string]afero.Fs{
		"file": c.FileFs,
		"http": c.HttpFs,
	}
}

func (c *SrcFs) NewFsManager() *fsm.Manager {
	m := fsm.NewManager()
	m.DefaultFsName = "file"
	m.AddFs("file", c.FileFs)
	m.AddFs("http", c.HttpFs)
	return m
}
