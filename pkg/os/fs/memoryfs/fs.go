package memoryfs

import (
	"github.com/spf13/afero"
)

const FsName = "memory"

type Fs struct {
	cfg *Config
	afero.MemMapFs
}

func NewFs(cfg *Config) (afero.Fs, error) {
	return &Fs{cfg: cfg}, nil
}

func (f *Fs) Name() string {
	return FsName
}

func (f *Fs) Tags() []string {
	return f.cfg.Tags
}

func Name() string {
	return FsName
}
