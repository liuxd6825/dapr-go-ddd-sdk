package localfs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/spf13/afero"
	"os"
)

const FsName = "local"

type Fs struct {
	afero.Fs
	cfg Config
}

func NewFs(cfg Config) (afero.Fs, error) {
	fs := afero.NewOsFs()
	return &Fs{Fs: fs, cfg: cfg}, nil
}

func (f *Fs) Name() string {
	return FsName
}

func (f *Fs) BasePath() string {
	return f.cfg.Path
}

// Open opens a file, returning it or an error, if any happens.
func (f *Fs) Open(name string) (afero.File, error) {
	file, err := f.Fs.Open(f.fileName(name))
	return file, err
}

// OpenFile opens a file using the given flags and the given mode.
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return f.Fs.OpenFile(f.fileName(name), flag, perm)
}

func (f *Fs) fileName(name string) string {
	fileName := fileutils.AbsPath(f.BasePath(), name)
	return fileName
}

func Name() string {
	return FsName
}
