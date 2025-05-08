package giteafs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"io"
	"os"
)

type file struct {
	fs       *Fs
	name     string
	data     []byte
	closed   bool
	fhOffset int64 // File handle specific offset
	flag     int
	perm     os.FileMode
	reader   io.Reader
}

func newFile(fs *Fs, name string, flag int, perm os.FileMode) *file {
	f := &file{fs: fs, name: name, closed: false, fhOffset: 0, flag: flag, perm: perm}
	return f
}

func (o *file) Close() error {
	if o.closed {
		return fsopts.ErrFileClosed
	}
	o.closed = true
	if o.flag == os.O_CREATE {
		return create(o.fs, o.name, o.data, "")
	} else if o.flag == os.O_WRONLY {
		return writeFile(o.fs, o.name, o.data, "")
	}
	return nil
}

func (o *file) Read(p []byte) (int, error) {
	return o.ReadAt(p, o.fhOffset)
}

func (o *file) ReadAt(p []byte, off int64) (n int, err error) {
	read, err := o.readAt(p, off)
	if err != nil {
		return 0, err
	}
	o.fhOffset += int64(read)
	return read, err
}

func (o *file) read(p []byte, off int64) (n int, err error) {
	if o.reader == nil {
		owner := o.fs.cfg.User
		repo := o.fs.cfg.Repo
		ref := o.fs.cfg.Branch
		path := o.name
		reader, _, err := o.fs.gitea.GetFileReader(owner, repo, path, ref)
		if err != nil {
			return 0, err
		}
		o.reader = reader
	}
	n, err = o.reader.Read(p)
	return n, err
}

func (o *file) readAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, os.ErrInvalid
	}
	if o.fhOffset == 0 {
		owner := o.fs.cfg.User
		repo := o.fs.cfg.Repo
		ref := o.fs.cfg.Branch
		path := o.name
		data, _, err := o.fs.gitea.GetFile(owner, repo, ref, path)
		if err != nil {
			return 0, err
		}
		o.data = data
		o.closed = true
	}
	readLen := int64(len(p))
	dataLen := int64(len(o.data))
	if off+readLen > dataLen {
		readLen = dataLen
	} else {
		readLen = off + readLen
	}
	if off >= dataLen {
		o.closed = true
		o.fhOffset = dataLen
		return 0, io.EOF
	}
	copy(p, o.data[off:readLen])
	return len(p), nil
}

func (o *file) Seek(newOffset int64, whence int) (int64, error) {
	if o.closed {
		return 0, fsopts.ErrFileClosed
	}

	if (whence == 0 && newOffset == o.fhOffset) || (whence == 1 && newOffset == 0) {
		return o.fhOffset, nil
	}

	err := o.Sync()
	if err != nil {
		return 0, err
	}
	stat, err := o.Stat()
	if err != nil {
		return 0, nil
	}

	switch whence {
	case 0:
		o.fhOffset = newOffset
	case 1:
		o.fhOffset += newOffset
	case 2:
		o.fhOffset = stat.Size() + newOffset
	}
	return o.fhOffset, nil
}

func (o *file) Write(p []byte) (n int, err error) {
	return o.WriteAt(p, o.fhOffset)
}

func (o *file) WriteAt(p []byte, off int64) (n int, err error) {
	if o.flag == os.O_RDONLY {
		o.flag = os.O_WRONLY
	}
	if o.closed {
		return 0, fsopts.ErrFileClosed
	}
	written := len(p)
	o.data = append(o.data, p...)
	o.fhOffset += int64(written)
	return written, err
}

func (o *file) Name() string {
	return o.name
}

func (o *file) Readdir(count int) ([]os.FileInfo, error) {
	return readDir(o.fs, o.name, count)
}

func (o *file) Readdirnames(count int) ([]string, error) {
	fileInfos, err := readDir(o.fs, o.name, count)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(fileInfos))
	for i, fileInfo := range fileInfos {
		names[i] = fileInfo.Name()
	}
	return names, nil
}

func (o *file) Stat() (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (o *file) Sync() error {
	return nil
}

func (o *file) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *file) WriteString(s string) (ret int, err error) {
	n, err := o.Write([]byte(s))
	return n, err
}
