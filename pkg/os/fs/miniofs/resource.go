package miniofs

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/minio/minio-go/v7"
)

func fullPath(f *Fs, name string) string {
	return strings.TrimLeft(f.rootPath+name, "/")
}

func objectInfoToFileInfo(o minio.ObjectInfo) fs.FileInfo {
	// Strip a trailing slash before slicing so directory markers yield the
	// directory's own name (e.g. "sub") instead of an empty string.
	key := strings.TrimRight(o.Key, "/")
	name := key
	if i := strings.LastIndex(key, "/"); i >= 0 {
		name = key[i+1:]
	}
	if name == "" {
		name = o.Key
	}
	// MinIO ObjectInfo does not carry a directory flag; we infer from the
	// key's trailing slash (zero-byte marker convention).
	isDir := strings.HasSuffix(o.Key, "/")
	mode := fs.FileMode(0644)
	if isDir {
		mode = fs.FileMode(os.ModeDir | 0755)
	}
	return fsopts.NewFileInfo(
		name,
		o.Size,
		mode,
		o.LastModified,
		isDir,
		nil,
	)
}

func getInfo(f *Fs, name string) (minio.ObjectInfo, error) {
	full := fullPath(f, name)
	if full == "" {
		return minio.ObjectInfo{}, errors.New("miniofs: empty object name")
	}
	stat, err := f.client.StatObject(context.Background(), f.cfg.Bucket, full, minio.GetObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return minio.ObjectInfo{}, fsopts.ErrObjectDoesNotExist
		}
		return minio.ObjectInfo{}, err
	}
	return stat, nil
}

func readFile(f *Fs, name string) ([]byte, error) {
	full := fullPath(f, name)
	if full == "" {
		return nil, fsopts.ErrEmptyObjectName
	}
	obj, err := f.client.GetObject(context.Background(), f.cfg.Bucket, full, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

func writeFile(f *Fs, name string, data []byte) error {
	full := fullPath(f, name)
	if full == "" {
		return fsopts.ErrEmptyObjectName
	}
	_, err := f.client.PutObject(
		context.Background(),
		f.cfg.Bucket,
		full,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
	)
	return err
}

func deleteFile(f *Fs, name string) error {
	full := fullPath(f, name)
	if full == "" {
		return fsopts.ErrEmptyObjectName
	}
	return f.client.RemoveObject(context.Background(), f.cfg.Bucket, full, minio.RemoveObjectOptions{})
}

// dirKey normalizes a path to a directory marker key ending with "/". An
// empty input yields an empty string.
func dirKey(s string) string {
	if s == "" {
		return ""
	}
	return strings.TrimRight(s, "/") + "/"
}

// putMarker uploads a zero-byte object whose key ends with "/" to mark a
// directory. Marker presence is how Readdir/Stat detect directories. The
// call is idempotent: MinIO accepts overwrites of an existing zero-byte
// object without complaint.
func putMarker(f *Fs, key string) error {
	if key == "" || !strings.HasSuffix(key, "/") {
		return errors.New("miniofs: marker key must be non-empty and end with /")
	}
	_, err := f.client.PutObject(
		context.Background(),
		f.cfg.Bucket,
		key,
		bytes.NewReader(nil),
		0,
		minio.PutObjectOptions{ContentType: "application/x-directory"},
	)
	return err
}
// deleteAll removes every object whose key starts with prefix, including
// any directory marker at the prefix itself. If prefix is empty the call
// returns an error because deleting the entire bucket is not what RemoveAll
// is supposed to do.
func deleteAll(f *Fs, prefix string) error {
	prefix = fullPath(f, prefix)
	if prefix == "" {
		return errors.New("miniofs: refusing to removeAll on empty prefix")
	}
	objectsCh := make(chan minio.ObjectInfo, 1000)
	go func() {
		defer close(objectsCh)
		for object := range f.client.ListObjects(context.Background(), f.cfg.Bucket, minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			objectsCh <- object
		}
	}()
	var firstErr error
	for rErr := range f.client.RemoveObjects(context.Background(), f.cfg.Bucket, objectsCh, minio.RemoveObjectsOptions{}) {
		if rErr.Err != nil {
			firstErr = rErr.Err
			break
		}
	}
	if firstErr != nil {
		return firstErr
	}
	return nil
}

func renameFile(f *Fs, oldName, newName string) error {
	src := fullPath(f, oldName)
	dst := fullPath(f, newName)
	if src == "" || dst == "" {
		return fsopts.ErrEmptyObjectName
	}
	_, err := f.client.CopyObject(context.Background(),
		minio.CopyDestOptions{Bucket: f.cfg.Bucket, Object: dst},
		minio.CopySrcOptions{Bucket: f.cfg.Bucket, Object: src},
	)
	if err != nil {
		return err
	}
	return f.client.RemoveObject(context.Background(), f.cfg.Bucket, src, minio.RemoveObjectOptions{})
}

func readDir(f *Fs, p string, count int) ([]os.FileInfo, error) {
	prefix := dirKey(fullPath(f, p))
	items := make([]os.FileInfo, 0)
	for object := range f.client.ListObjects(context.Background(), f.cfg.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: false,
	}) {
		// Skip the directory's own marker: os.File.Readdir never returns
		// the directory itself, only its children.
		if object.Key == prefix {
			continue
		}
		items = append(items, objectInfoToFileInfo(object))
	}
	// os.File.Readdir returns entries in lexical (byte-wise) order on Linux.
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})
	if count > 0 && len(items) > count {
		items = items[:count]
	}
	return items, nil
}

func detectContentType(name string) string {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".html", ".htm":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".xml":
		return "application/xml"
	case ".txt":
		return "text/plain"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".zip":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// writeStream starts a streaming PutObject against the given object name and
// returns the writer end of the pipe. The caller is expected to Close the
// returned PipeWriter once all data has been written. errCh will receive
// exactly one value (nil or non-nil) when the PutObject goroutine finishes.
func writeStream(f *Fs, name string) (io.WriteCloser, <-chan error, error) {
	full := fullPath(f, name)
	if full == "" {
		return nil, nil, fsopts.ErrEmptyObjectName
	}
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)
	go func() {
		_, err := f.client.PutObject(
			context.Background(),
			f.cfg.Bucket,
			full,
			pr,
			-1,
			minio.PutObjectOptions{ContentType: detectContentType(name)},
		)
		errCh <- err
	}()
	return pw, errCh, nil
}

// errCode unwraps a minio ErrorResponse to a stable code, or "" if err is nil
// or not a minio error.
func errCode(err error) string {
	if err == nil {
		return ""
	}
	return minio.ToErrorResponse(err).Code
}

var _ = errCode
