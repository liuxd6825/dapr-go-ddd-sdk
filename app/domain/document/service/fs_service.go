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

const DocumentFsName = "documentIoStore"

var _fsService *FsService
var _fsServiceOnce sync.Once

func NewFsService() *FsService {
	_fsServiceOnce.Do(func() {
		_fsService = newFsService()
	})
	return _fsService
}

func newFsService() *FsService {
	fsService := &FsService{}
	fsService.Init()
	return fsService
}

func (s *FsService) Init() *FsService {
	docFs, err := fs_pkg.NewFsPkg(env.GetEnv(), DocumentFsName)
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

func (s *FsService) RemoveFile(filename string) {
	s.docFs.RemoveFile(filename)
}
func (s *FsService) RemoveAll(name string) {
	s.docFs.RemoveAll(name)
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

func (s *FsService) MoveDir(source string, target string) error {
	return s.docFs.MoveDir(source, target)
}

// Download 下载文件取Web  objectName:带路径的文件完整名称  saveFileName:保存到本地时的文件名称
func (s *FsService) Download(ictx iris.Context, fullFileName string, saveFileName string) error {
	readFile, err := s.docFs.Open(fullFileName, os.O_RDONLY, 0666)

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
