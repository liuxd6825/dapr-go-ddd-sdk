package service

import (
	_ "embed"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/fs_pkg"
)

type FsService struct {
	portalFs fspkg.IFsPkg
}

var _fsService *FsService
var _fsServiceOnce sync.Once

const FilePerm = os.ModePerm

func NewFsService(dbKey string) *FsService {
	_fsServiceOnce.Do(func() {
		_fsService = newFsService(dbKey)
	})
	return _fsService
}

func newFsService(dbKey string) *FsService {
	fsService := &FsService{}
	fsService.Init()
	return fsService
}

func (s *FsService) Init() *FsService {
	portalFs, err := fs_pkg.NewFsPkg(envs.GetEnv(), "portalIoStore")
	if err != nil {
		panic("portalIoStore fs not exist in config")
	}
	s.portalFs = portalFs
	return s
}

func (s *FsService) Create(fileName string) afero.File {
	pathName := filepath.Dir(fileName)
	exists := s.Exists(pathName)
	if !exists && pathName != "." {
		s.MkdirAll(pathName)
	}
	return s.portalFs.Create(fileName)
}

func (s *FsService) MkdirAll(path string) {
	s.portalFs.MkdirAll(path, FilePerm)
}

func (s *FsService) RemoveFile(filename string) {
	s.portalFs.RemoveFile(filename)
}

func (s *FsService) WriteFile(fileName string, data any) error {
	s.portalFs.WriteFile(fileName, data)
	return nil
}

func (s *FsService) Exists(fileName string) bool {
	return s.portalFs.Exists(fileName)
}

func (s *FsService) Download(ictx iris.Context, fullFileName string, saveFileName string) error {
	readFile, err := s.portalFs.Open("/"+fullFileName, os.O_RDONLY, FilePerm)
	if err != nil {
		return err
	}
	defer func() {
		_ = readFile.Close()
	}()
	fileInfo, err := readFile.Stat()
	if err != nil {
		return err
	}
	ictx.ResponseWriter().Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(saveFileName))
	ictx.ServeContentWithRate(readFile, filepath.Base(saveFileName), fileInfo.ModTime(), 1024*iris.KB, 1024*iris.KB)
	return nil
}
