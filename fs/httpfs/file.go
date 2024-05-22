package httpfs

import (
	"github.com/spf13/afero"
	"os"
)

type file struct {
	fs   *Fs
	name string
}

func newFile(fs *Fs, name string) afero.File {
	return &file{fs: fs, name: name}
}

func (f *file) Close() error {
	return nil
}

func (f *file) Read(p []byte) (n int, err error) {
	data, err := f.fs.readFile(f.name)
	if err != nil {
		return 0, err
	}
	p = data
	return len(data), nil
}

func (f *file) ReadAt(p []byte, off int64) (n int, err error) {
	return f.Read(p)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (f *file) Write(p []byte) (n int, err error) {
	err = f.fs.writeFile(f.name, string(p), "")
	n = len(p)
	return
}

func (f *file) WriteAt(p []byte, off int64) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *file) Readdirnames(n int) ([]string, error) {
	return nil, nil
}

func (f *file) Stat() (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *file) Sync() error {
	return nil
}

func (f *file) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *file) WriteString(s string) (ret int, err error) {
	//TODO implement me
	panic("implement me")
}
