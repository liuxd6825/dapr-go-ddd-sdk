package localfs

import "github.com/spf13/afero"

func NewFs(cfg *Config) (afero.Fs, error) {
	return afero.NewBasePathFs(afero.NewOsFs(), cfg.Path), nil
}
