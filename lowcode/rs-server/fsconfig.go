package rs_server

import "github.com/spf13/afero"

type FsConfig struct {
	FileFs afero.Fs
	HttpFs afero.Fs
}

func NewFsConfig(fileFs, httpFs afero.Fs) *FsConfig {
	return &FsConfig{FileFs: fileFs, HttpFs: httpFs}
}

func (c *FsConfig) ToMap() map[string]afero.Fs {
	return map[string]afero.Fs{
		"file": c.FileFs,
		"http": c.HttpFs,
	}
}
