package miniofs

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/afero"
)

// Fs is an afero.Fs backed by a MinIO/S3-compatible object store. The
// minio.Client is supplied by the caller so the same instance can be shared
// between filesystem and other MinIO consumers in the process.
type Fs struct {
	cfg      Config
	rootPath string
	client   *minio.Client
	minio    *env.Minio
	env      *env.Env
}

func NewFs(env *env.Env, minioKey string, bucketKey string) (*Fs, error) {
	m, ok := env.GetMinioByKey(minioKey)
	if !ok {
		return nil, errors.New(fmt.Sprintf("在配置文件中没有定义 minio.%s 节.", minioKey))
	}
	if m.Client == nil {
		return nil, errors.New(fmt.Sprintf(""))
	}
	bucketName, ok := m.Buckets[bucketKey]
	if !ok {
		return nil, errors.New(fmt.Sprintf("在配置文件中没有定义 minio.%s.buckets.%s 值.", minioKey, bucketKey))
	}

	cfg := Config{
		Bucket: bucketName,
	}

	if m.Client == nil {
		return nil, errors.New("miniofs.NewFs: client is nil")
	}
	return &Fs{cfg: cfg, rootPath: cfg.RootPath, client: m.Client, minio: m, env: env}, nil
}

// Compile-time assertion that *Fs implements afero.Fs.
var _ afero.Fs = (*Fs)(nil)

func (f *Fs) Tags() []string {
	return f.cfg.Tags
}

func (f *Fs) GetS3FileURL(fileName string) string {
	/*
		CALL apoc.load.parquet('s3://minioadmin:minioadmin@192.168.120.224:9000/import-data/record/record_hUe2yxR81YNB89AY3sWRIM4U.parquet')
		YIELD value
		RETURN value LIMIT 100;
	*/
	return fmt.Sprintf(
		"s3://%s:%s@%s/%s/%s",
		f.minio.AccessKey,
		f.minio.SecretKey,
		f.minio.Endpoint,
		f.cfg.Bucket,
		fileName,
	)
}

func (f *Fs) GetFsType() string {
	return "minio"
}

func (f *Fs) GetBucketName() string {
	return f.cfg.Bucket
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

// Create creates or truncates the named file and opens it for writing.
func (f *Fs) Create(name string) (afero.File, error) {
	return newFile(f, name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
}

// Mkdir creates a zero-byte marker object whose key ends with "/" so the
// directory can be detected by subsequent Stat and Readdir calls. The
// behaviour mirrors os.Mkdir: a single directory is created and the call
// succeeds when the marker already exists.
func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	return putMarker(f, dirKey(fullPath(f, name)))
}

// MkdirAll creates a zero-byte marker for every level of path. Mirrors
// os.MkdirAll: intermediate directories are also created. An empty path is
// a no-op so callers can pass a path that may already be rooted.
func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	full := fullPath(f, path)
	if full == "" {
		return nil
	}
	parts := strings.Split(strings.Trim(full, "/"), "/")
	cur := ""
	for _, p := range parts {
		cur += p + "/"
		if err := putMarker(f, cur); err != nil {
			return err
		}
	}
	return nil
}

// Open opens the named file for reading.
func (f *Fs) Open(name string) (afero.File, error) {
	return newFile(f, name, os.O_RDONLY, 0)
}

// OpenFile dispatches to newFile, which picks the streaming PutObject or
// GetObject path based on the flag bits.
func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return newFile(f, name, flag, perm)
}

// Remove deletes a single object. It does not recurse; use RemoveAll for
// directory-like removals.
func (f *Fs) Remove(name string) error {
	return deleteFile(f, name)
}

// RemoveAll deletes every object whose key starts with name. Because MinIO has
// no real directories, the prefix itself is also removed when it happens to
// be a zero-byte marker.
func (f *Fs) RemoveAll(path string) error {
	return deleteAll(f, path)
}

// Rename copies the object to the new name and then removes the original.
func (f *Fs) Rename(oldname, newname string) error {
	return renameFile(f, oldname, newname)
}

// Stat returns the metadata of the named object. If the bare name does not
// resolve, the call retries with a trailing "/" so directory markers created
// by Mkdir/MkdirAll are visible to callers that follow the Go os.Stat
// convention of "name without trailing slash".
func (f *Fs) Stat(name string) (os.FileInfo, error) {
	info, err := getInfo(f, name)
	if err != nil && !errors.Is(err, fsopts.ErrObjectDoesNotExist) {
		return nil, err
	}
	if err != nil {
		// Retry as a directory marker (key ending with "/").
		info, err = getInfo(f, dirKey(fullPath(f, name)))
		if err != nil {
			return nil, err
		}
	}
	return objectInfoToFileInfo(info), nil
}

// Chmod is a no-op: object stores do not have POSIX permissions.
func (f *Fs) Chmod(name string, mode os.FileMode) error {
	return nil
}

// Chown is a no-op: object stores do not have POSIX ownership.
func (f *Fs) Chown(name string, uid, gid int) error {
	return nil
}

// Chtimes is a no-op: object stores do not expose arbitrary access/modification
// time writes through the basic API.
func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}

func Name() string {
	return "minio"
}
