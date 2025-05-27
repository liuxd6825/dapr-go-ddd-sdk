package service

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
)

type Options struct {
	WorkPath string // 当前工作目录
}

type IFileService interface {
	CreateFile(ctx context.Context, fileName string, data []byte, opts ...*Options) error
	ReadFile(ctx context.Context, fileName string, opts ...*Options) ([]byte, error)
	WriteFile(ctx context.Context, fileName string, data []byte, opts ...*Options) error
	WriteJson(ctx context.Context, fileName string, data []byte, opts ...*Options) error
	RemoveFile(ctx context.Context, fileName string, opts ...*Options) error
	RenameFile(ctx context.Context, fileName string, newName string, opts ...*Options) error
	CreatePath(ctx context.Context, folderName string, opts ...*Options) error
	ReadPath(ctx context.Context, folderName string, isDeep bool, opts ...*Options) []*fspkg.FileInfo
	RemovePath(ctx context.Context, folderName string, isDeep bool, opts ...*Options) error
}

type fileService struct {
	env      *env.Env
	rootPath string
	fsPkg    fspkg.IFsPkg
}

func NewFileService(env *env.Env, rootPath string) IFileService {
	fsPkg, err := fspkg.NewFsPkg(env, rootPath)
	if err != nil {
		panic(err)
	}

	service := &fileService{
		env:      env,
		rootPath: rootPath,
		fsPkg:    fsPkg,
	}
	return service
}

func (s *fileService) CreateFile(ctx context.Context, fileName string, data []byte, opts ...*Options) error {
	fs := s.fsPkg.Create(fileName, newFsOptions(opts...))
	_, err := fs.Write(data)
	return err
}

func (s *fileService) ReadFile(ctx context.Context, fileName string, opts ...*Options) ([]byte, error) {
	return s.fsPkg.ReadFile(fileName, newFsOptions(opts...)), nil
}

func (s *fileService) WriteFile(ctx context.Context, fileName string, schema []byte, opts ...*Options) error {
	fs := s.fsPkg.Create(fileName, newFsOptions(opts...))
	_, err := fs.Write(schema)
	return err
}

func (s *fileService) WriteJson(ctx context.Context, fileName string, data []byte, opts ...*Options) error {
	s.fsPkg.WriteJson(fileName, data, newFsOptions(opts...))
	return nil
}

func (s *fileService) RemoveFile(ctx context.Context, fileName string, opts ...*Options) error {
	s.fsPkg.RemoveFile(fileName, newFsOptions(opts...))
	return nil
}

func (s *fileService) RenameFile(ctx context.Context, fileName string, newName string, opts ...*Options) error {
	return s.fsPkg.Rename(fileName, newName, newFsOptions(opts...))
}

func (s *fileService) CreatePath(ctx context.Context, pathName string, opts ...*Options) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	s.fsPkg.Mkdir(pathName, 1, newFsOptions(opts...))
	return nil
}

func (s *fileService) ReadPath(ctx context.Context, pathName string, isDeep bool, opts ...*Options) []*fspkg.FileInfo {
	if isDeep {
		return s.fsPkg.ReadAllPath(pathName, newFsOptions(opts...))
	}
	return s.fsPkg.ReadPath(pathName, newFsOptions(opts...))
}

func (s *fileService) RemovePath(ctx context.Context, pathName string, isDeep bool, opts ...*Options) error {
	s.fsPkg.RemoveAll(pathName, newFsOptions(opts...))
	return nil
}

func newFsOptions(opts ...*Options) *fsopts.Options {
	res := fsopts.NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		res.WorkPath = opt.WorkPath
	}
	return res
}
