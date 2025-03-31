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
	ReadPath(ctx context.Context, pathName string) []*fspkg.FileInfo
	ReadAllPath(ctx context.Context, pathName string) []*fspkg.FileInfo
	GetProperties(ctx context.Context, fileName string) map[string]*jsonschema.Schema
	GetSchema(ctx context.Context, fileName string) *jsonschema.Schema
}

type schemaService struct {
	fileService IFileService
}

func NewSchemaService(env env.IEnvConfig, rootPath string) ISchemaService {
	return &schemaService{
		fileService: NewFileService(env, rootPath),
	}
}

func (s *schemaService) ReadFile(ctx context.Context, fileName string) ([]byte, error) {
	fileName = "/" + fileName
	data, err := s.fileService.ReadFile(ctx, fileName)
	return data, err
}

func (s *schemaService) WriteFile(ctx context.Context, fileName string, data []byte) error {
	fileName = "/" + fileName
	return s.fileService.WriteFile(ctx, fileName, data)
}

func (s *schemaService) ReadPath(ctx context.Context, pathName string) []*fspkg.FileInfo {
	pathName = "/" + pathName
	return s.fileService.ReadPath(ctx, pathName, false)
}

func (s *schemaService) ReadAllPath(ctx context.Context, pathName string) []*fspkg.FileInfo {
	return s.fileService.ReadPath(ctx, pathName, true)
}

func (s *schemaService) GetProperties(ctx context.Context, fileName string) map[string]*jsonschema.Schema {
	sch := s.GetSchema(ctx, fileName)
	return sch.GetAllProperties()
}

func (s *schemaService) GetSchema(ctx context.Context, fileName string) *jsonschema.Schema {
	fileName = "/" + fileName
	data, err := s.fileService.ReadFile(ctx, fileName)
	if err != nil {
		panic(err)
	}
	sch := schema.NewJsonSchemaWithBytes(fileName, data)
	return sch
}
