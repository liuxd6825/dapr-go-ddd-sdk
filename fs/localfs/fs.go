package localfs

import "github.com/spf13/afero"

const FsName = "local"

type Fs struct {
	afero.Fs
}

func NewFs(cfg *Config) (afero.Fs, error) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), cfg.Path)
	return &Fs{fs}, nil
}

func (f *Fs) Name() string {
	return FsName
}

func Name() string {
	return FsName
}
