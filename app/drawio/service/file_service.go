package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"os"
	"path/filepath"
	"sync"
)

type FileService struct {
	drawFs fspkg.IFsPkg
}

var _fileService *FileService
var _fileServiceOnce sync.Once

func NewFileService() *FileService {
	_fileServiceOnce.Do(func() {
		_fileService = newFileService()
	})
	return _fileService
}

func newFileService() *FileService {
	fileService := &FileService{}
	fileService.Init()
	return fileService
}

func (s *FileService) Init() *FileService {
	drawFs, err := fs_pkg.NewFsPkg(env.GetEnv(), "drawIoStore")
	if err != nil {
		panic("drawIoStore fs not exist in config")
	}
	s.drawFs = drawFs
	return s
}

func (s *FileService) Save(fileName string, content string) error {
	if fileName == "" {
		return errors.New("file name is empty")
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".drawio"
	}
	pathName := filepath.Dir(fileName)
	exists := s.drawFs.Exists(pathName)
	if !exists && pathName != "." {
		s.drawFs.Mkdir(pathName, os.ModePerm)
	}
	s.drawFs.WriteFile(fileName, content)
	return nil
}

func (s *FileService) Read(fileName string) (content []byte, err error) {
	if fileName == "" {
		return nil, errors.New("file name is empty")
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".drawio"
	}
	content = s.drawFs.ReadFile(fileName)
	return content, err
}
