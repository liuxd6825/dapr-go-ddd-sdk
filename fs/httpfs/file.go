package httpfs

import (
	"io"
	"os"
)

type File struct {
	name     string
	fs       *Fs
	data     []byte
	fhOffset int64
}

func newFile(fs *Fs, name string) *File {
	return &File{name: name, fs: fs}
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Close() error {
	return nil
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.readAt(p, f.fhOffset)
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return f.readAt(p, f.fhOffset)
}

func (f *File) readAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, os.ErrInvalid
	}
	if f.fhOffset == 0 {
		f.data, err = readFile(f.fs, f.name)
	}
	readLen := int64(len(p))
	dataLen := int64(len(f.data))
	if off+readLen > dataLen {
		readLen = dataLen
	} else {
		readLen = off + readLen
	}
	if off >= dataLen {
		f.fhOffset = dataLen
		return 0, io.EOF
	}
	copy(p, f.data[off:readLen])
	return len(p), nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) WriteAt(p []byte, off int64) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Readdirnames(n int) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Stat() (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *File) Sync() error {
	//TODO implement me
	panic("implement me")
}

func (f *File) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (f *File) WriteString(s string) (ret int, err error) {
	//TODO implement me
	panic("implement me")
}
