package rs_server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"github.com/spf13/afero"
)

// SrcFsConfig
// @Description: source code file system config
type SrcFsConfig struct {
	FileFs afero.Fs
	HttpFs afero.Fs
}

func NewSrcFsConfig(fileFs, httpFs afero.Fs) *SrcFsConfig {
	return &SrcFsConfig{FileFs: fileFs, HttpFs: httpFs}
}

func (c *SrcFsConfig) ToMap() map[string]afero.Fs {
	return map[string]afero.Fs{
		"file": c.FileFs,
		"http": c.HttpFs,
	}
}

func (c *SrcFsConfig) NewFsManager() *fsm.Manager {
	m := &fsm.Manager{}
	m.Add("file", c.FileFs)
	m.Add("http", c.HttpFs)
	return m
}
