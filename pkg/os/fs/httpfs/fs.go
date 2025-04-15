package httpfs

import (
	"github.com/go-resty/resty/v2"
	"github.com/spf13/afero"
	"os"
	"time"
)

type Fs struct {
	cfg    *Config
	client *resty.Client
}

func NewFs(cfg *Config) (afero.Fs, error) {
	return &Fs{cfg: cfg, client: resty.New()}, nil
}

func (f *Fs) Tags() []string {
	return f.cfg.Tags
}

func (f *Fs) GetRootPath() string {
	return f.cfg.BaseUrl
}

func (f *Fs) absFile(path string) string {
	return f.cfg.BaseUrl + path
}

func (f *Fs) Create(name string) (afero.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Open(name string) (afero.File, error) {
	return newFile(f, name), nil
}

func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) RemoveAll(path string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Rename(oldname, newname string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Stat(name string) (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Name() string {
	return "http"
}

func (f *Fs) Chmod(name string, mode os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Chown(name string, uid, gid int) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	//TODO implement me
	panic("implement me")
}
