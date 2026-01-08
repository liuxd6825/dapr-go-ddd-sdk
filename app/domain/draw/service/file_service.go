package service

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/hserver/pkg/fs_pkg"

	"os"
	"path/filepath"
	"sync"
)

type FileService struct {
	drawFs fspkg.IFsPkg
}

//go:embed resource/new.drawio
var _newContent string

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

func (s *FileService) Create(ctx context.Context, caseId, drawId, fileName string) error {
	return s.Save(ctx, caseId, fileName, drawId, _newContent)
}

func (s *FileService) Save(ctx context.Context, caseId, drawId string, fileName string, content string) error {
	if fileName == "" {
		return errors.New("file name is empty")
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".drawio"
	}
	tenantId, _ := appctx.GetTenantId(ctx)
	fileName = fmt.Sprintf("/%s/%s/%s.drawio", tenantId, caseId, drawId)
	pathName := filepath.Dir(fileName)
	exists := s.drawFs.Exists(pathName)
	if !exists && pathName != "." {
		s.drawFs.MkdirAll(pathName, os.ModePerm)
	}
	s.drawFs.WriteFile(fileName, content)
	return nil
}

func (s *FileService) Read(ctx context.Context, caseId string, fileName string) (content []byte, err error) {
	if fileName == "" {
		return nil, errors.New("file name is empty")
	}
	if filepath.Ext(fileName) == "" {
		fileName += ".drawio"
	}
	tenantId, _ := appctx.GetTenantId(ctx)
	fileName = fmt.Sprintf("/%s/%s/%s", tenantId, caseId, fileName)
	content = s.drawFs.ReadFile(fileName)
	return content, err
}
