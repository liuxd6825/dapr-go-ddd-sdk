package localfs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils/mapstructure"
	"github.com/spf13/afero"
	"os"
)

type Config struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type Fs struct {
	*afero.OsFs
	cfg      Config
	rootPath string
}

func NewConfig(metadata map[string]any) (*Config, error) {
	cfg := &Config{}
	err := mapstructure.Decode(metadata, cfg)
	if err != nil {
		return nil, err
	}
	vErr := errors.NewVerifyError()
	vErr.Message = fmt.Sprintf("localfs.config: %s", cfg.Name)
	if cfg.Name == "" {
		vErr.AppendField("id", "missing id")
	}
	if cfg.Path == "" {
		vErr.AppendField("path", "missing url")
	}
	if vErr.HasError() {
		return nil, vErr
	}
	return cfg, err
}

func NewFs(cfg Config) (afero.Fs, error) {
	osFs, ok := afero.NewOsFs().(*afero.OsFs)
	if !ok {
		return nil, fmt.Errorf("localfs is not a OsFs")
	}
	return &Fs{OsFs: osFs, cfg: cfg, rootPath: cfg.Path}, nil
}

func (f *Fs) NewWatcher() (fs.Watcher, error) {
	return NewFileWatcher(f), nil
}

func (f *Fs) GetFsType() string {
	return "local"
}

func (f *Fs) Name() string {
	return f.GetFsType()
}

func (f *Fs) GetRootPath() string {
	return f.rootPath
}

func (f *Fs) SetRootPath(path string) {
	f.rootPath = path
}

// Open opens a file, returning it or an error, if any happens.
func (f *Fs) Open(filename string) (afero.File, error) {
	file, err := f.OsFs.Open(f.rootPath + filename)
	return file, err
}

func (f *Fs) Create(name string) (afero.File, error) {
	return f.OsFs.Create(f.rootPath + name)
}

func (f *Fs) Mkdir(filename string, perm os.FileMode) error {
	return f.OsFs.Mkdir(f.rootPath+filename, perm)
}

func (f *Fs) RemoveAll(filename string) error {
	return f.OsFs.RemoveAll(f.rootPath + filename)
}

func (f *Fs) Remove(filename string) error {
	return f.OsFs.Remove(f.rootPath + filename)
}

func (f *Fs) MkdirAll(filename string, mode os.FileMode) error {
	return f.OsFs.MkdirAll(f.rootPath+filename, mode)
}

func (f *Fs) Rename(oldName, newName string) error {
	return f.OsFs.Rename(f.rootPath+oldName, f.rootPath+newName)
}

func (f *Fs) Stat(name string) (os.FileInfo, error) {
	return f.OsFs.Stat(f.rootPath + name)
}

func (f *Fs) OpenFile(filename string, flag int, perm os.FileMode) (afero.File, error) {
	return f.OsFs.OpenFile(f.rootPath+filename, flag, perm)
}

func (f *Fs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	return f.OsFs.LstatIfPossible(f.rootPath + name)
}

func Name() string {
	return "local"
}
