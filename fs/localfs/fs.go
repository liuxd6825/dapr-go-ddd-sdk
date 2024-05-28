package localfs

import (
	"github.com/spf13/afero"
)

const FsName = "local"

type Fs struct {
	afero.Fs
	cfg *Config
}

func NewFs(cfg *Config) (afero.Fs, error) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), cfg.Path)
	return &Fs{Fs: fs, cfg: cfg}, nil
}

func (f *Fs) Name() string {
	return FsName
}

func (f *Fs) BasePath() string {
	return f.cfg.Path
}

func Name() string {
	return FsName
}
