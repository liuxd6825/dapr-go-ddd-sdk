package miniofs

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/minio/minio-go/v7"
)

type file struct {
	fs     *Fs
	name   string
	flag   int
	perm   os.FileMode
	closed bool

	// reader is the *minio.Object returned by GetObject; non-nil for read paths.
	reader *minio.Object
	// writer is the io.PipeWriter returned by writeStream; non-nil for write paths.
	writer io.WriteCloser
	// errCh receives the result of the streaming PutObject goroutine; closed by Close.
	errCh <-chan error

	// offset tracks the read or write position. For a reader, it mirrors the
	// underlying *minio.Object's seek state. For a writer it equals the bytes
	// already pushed into the pipe.
	offset int64
}

func newFile(fsys *Fs, name string, flag int, perm os.FileMode) (*file, error) {
	f := &file{fs: fsys, name: name, flag: flag, perm: perm}

	writable := flag&(os.O_WRONLY|os.O_RDWR) != 0 || flag&(os.O_CREATE|os.O_TRUNC) != 0
	if writable {
		w, errCh, err := writeStream(fsys, name)
		if err != nil {
			return nil, err
		}
		f.writer = w
		f.errCh = errCh
		return f, nil
	}

	// Read path
	full := fullPath(fsys, name)
	if full == "" {
		return nil, fsopts.ErrEmptyObjectName
	}
	obj, err := fsys.client.GetObject(context.Background(), fsys.cfg.Bucket, full, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	f.reader = obj
	return f, nil
}

func (f *file) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true
	if f.writer != nil {
		// Closing the pipe signals EOF to the PutObject goroutine.
		if err := f.writer.Close(); err != nil {
			return err
		}
		if f.errCh != nil {
			if err, ok := <-f.errCh; ok && err != nil {
				return err
			}
		}
		return nil
	}
	if f.reader != nil {
		return f.reader.Close()
	}
	return nil
}

func (f *file) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fsopts.ErrFileClosed
	}
	if f.reader == nil {
		return 0, errors.New("miniofs: file not opened for reading")
	}
	n, err := f.reader.Read(p)
	f.offset += int64(n)
	return n, err
}

func (f *file) ReadAt(p []byte, off int64) (int, error) {
	if f.closed {
		return 0, fsopts.ErrFileClosed
	}
	if f.reader == nil {
		return 0, errors.New("miniofs: file not opened for reading")
	}
	return f.reader.ReadAt(p, off)
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	if f.closed {
		return 0, fsopts.ErrFileClosed
	}
	if f.reader == nil {
		return 0, errors.New("miniofs: Seek is only supported on read-only files")
	}
	n, err := f.reader.Seek(offset, whence)
	if err == nil {
		f.offset = n
	}
	return n, err
}

func (f *file) Write(p []byte) (int, error) {
	if f.closed {
		return 0, fsopts.ErrFileClosed
	}
	if f.writer == nil {
		return 0, errors.New("miniofs: file not opened for writing")
	}
	n, err := f.writer.Write(p)
	f.offset += int64(n)
	return n, err
}

func (f *file) WriteAt(p []byte, off int64) (int, error) {
	if f.closed {
		return 0, fsopts.ErrFileClosed
	}
	if f.writer == nil {
		return 0, errors.New("miniofs: file not opened for writing")
	}
	// Streaming uploads do not support random-access writes; require the
	// caller to write sequentially starting at offset 0.
	if off != f.offset {
		return 0, fsopts.ErrOutOfRange
	}
	return f.Write(p)
}

func (f *file) Name() string {
	return f.name
}

func (f *file) Readdir(count int) ([]os.FileInfo, error) {
	if f.reader == nil {
		return nil, errors.New("miniofs: Readdir requires a directory-style file")
	}
	return readDir(f.fs, f.name, count)
}

func (f *file) Readdirnames(n int) ([]string, error) {
	items, err := f.Readdir(n)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(items))
	for _, it := range items {
		names = append(names, it.Name())
	}
	return names, nil
}

func (f *file) Stat() (os.FileInfo, error) {
	if f.closed {
		return nil, fsopts.ErrFileClosed
	}
	if f.reader != nil {
		info, err := f.reader.Stat()
		if err != nil {
			return nil, err
		}
		return objectInfoToFileInfo(info), nil
	}
	// Writer path: we only know the bytes pushed so far through the pipe.
	return fsopts.NewFileInfo(
		f.name,
		f.offset,
		fs.FileMode(f.perm),
		time.Now(),
		false,
		nil,
	), nil
}

func (f *file) Sync() error {
	// PutObject uploads synchronously into the pipe; there is no
	// client-side buffer to flush.
	return nil
}

func (f *file) Truncate(size int64) error {
	// Object stores do not support partial truncation.
	return fsopts.ErrOutOfRange
}

func (f *file) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}
