package service

import (
	_ "embed"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type FsService struct {
	docFs fspkg.IFsPkg
}

var _fileService *FsService
var _fileServiceOnce sync.Once

func NewFsService() *FsService {
	_fileServiceOnce.Do(func() {
		_fileService = newFsService()
	})
	return _fileService
}

func newFsService() *FsService {
	fileService := &FsService{}
	fileService.Init()
	return fileService
}

func (s *FsService) Init() *FsService {
	docFs, err := fs_pkg.NewFsPkg(env.GetEnv(), "documentIoStore")
	if err != nil {
		panic("documentIoStore fs not exist in config")
	}
	s.docFs = docFs
	return s
}

func (s *FsService) Exists(fileName string) bool {
	return s.docFs.Exists(fileName)
}

func (s *FsService) Create(fileName string) afero.File {
	return s.docFs.Create(fileName)
}

func (s *FsService) Rename(oldName, newName string) error {
	return s.docFs.Rename(oldName, newName)
}

func (s *FsService) MkdirAll(path string) {
	s.docFs.MkdirAll(path, 0666)
}

func (s *FsService) WriteAt(fileName string, data []byte, chunkIndex string, chunkSize string) error {
	writeFile, err := s.docFs.Open(fileName, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	idx, _ := strconv.ParseInt(chunkIndex, 10, 64)
	size, _ := strconv.ParseInt(chunkSize, 10, 64)

	off := idx * size

	_, err = s.docFs.WriteAt(writeFile, data, off)
	if err != nil {
		return err
	}

	err = writeFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *FsService) Download(ictx iris.Context, objectName string, fileName string) error {
	writeFile, err := s.docFs.Open("/"+objectName, os.O_WRONLY, 0666)
	//writeFile, err := os.Open(fsopts.GetAbsPath("/" + objectName))
	if err != nil {
		return err
	}

	fileInfo, err := writeFile.Stat()

	if err != nil {
		return err
	}

	ictx.ResponseWriter().Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(fileName))

	ictx.ServeContentWithRate(writeFile, filepath.Base(fileName), fileInfo.ModTime(), 0, 0)

	err = writeFile.Close()
	if err != nil {
		return err
	}

	return nil
}
