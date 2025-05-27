package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
)

type ISchemaService interface {
	ReadFile(ctx context.Context, fileName string) ([]byte, error)
	WriteFile(ctx context.Context, fileName string, data []byte) error
	WriteJson(ctx context.Context, fileName string, data []byte) error
	ReadPath(ctx context.Context, pathName string) []*fspkg.FileInfo
	ReadAllPath(ctx context.Context, pathName string) []*fspkg.FileInfo
	GetSchema(ctx context.Context, fileName string) *jsonschema.Schema
	CreatePath(ctx context.Context, pathName string) error
	CreateFile(ctx context.Context, fileName string) error
	RenameFile(ctx context.Context, fileName string, newFileName string) error
	RemoveFile(ctx context.Context, fileName string) error
	RemovePath(ctx context.Context, pathName string) error
}

type schemaService struct {
	fileService IFileService
}

func NewSchemaService(env *env.Env, rootPath string) ISchemaService {
	return &schemaService{
		fileService: NewFileService(env, rootPath),
	}
}

func (s *schemaService) ReadFile(ctx context.Context, fileName string) ([]byte, error) {
	data, err := s.fileService.ReadFile(ctx, fileName)
	return data, err
}

func (s *schemaService) WriteFile(ctx context.Context, fileName string, data []byte) error {
	return s.fileService.WriteFile(ctx, fileName, data)
}

func (s *schemaService) WriteJson(ctx context.Context, fileName string, data []byte) error {
	return s.fileService.WriteJson(ctx, fileName, data)
}

func (s *schemaService) ReadPath(ctx context.Context, pathName string) []*fspkg.FileInfo {
	return s.fileService.ReadPath(ctx, pathName, false)
}

func (s *schemaService) ReadAllPath(ctx context.Context, pathName string) []*fspkg.FileInfo {
	return s.fileService.ReadPath(ctx, pathName, true)
}

func (s *schemaService) GetSchema(ctx context.Context, fileName string) *jsonschema.Schema {
	data, err := s.fileService.ReadFile(ctx, fileName)
	if err != nil {
		panic(err)
	}
	sch := schema.NewJsonSchemaWithBytes(fileName, data)
	return sch
}

func (s *schemaService) CreatePath(ctx context.Context, pathName string) error {
	return s.fileService.CreatePath(ctx, pathName)
}

func (s *schemaService) CreateFile(ctx context.Context, fileName string) error {
	return s.fileService.CreateFile(ctx, fileName, nil)
}

func (s *schemaService) RenameFile(ctx context.Context, fileName string, newFileName string) error {
	return s.fileService.RenameFile(ctx, fileName, newFileName)
}

func (s *schemaService) RemoveFile(ctx context.Context, fileName string) error {
	return s.fileService.RemoveFile(ctx, fileName)
}

func (s *schemaService) RemovePath(ctx context.Context, pathName string) error {
	return s.fileService.RemovePath(ctx, pathName, true)
}
