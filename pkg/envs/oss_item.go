package envs

import (
	"github.com/minio/minio-go/v7"
)

type OSSItem interface {
	GetOssType() OssType
	GetOssKey() string
	GetMinioClient() *minio.Client
	GetConfig() any
}

type ossItem struct {
	fsKey       string
	ossType     OssType
	minioClient *minio.Client
	config      any
}

func (f *ossItem) GetOssType() OssType {
	return f.ossType
}

func (f *ossItem) GetOssKey() string {
	return f.fsKey
}

func (f *ossItem) GetMinioClient() *minio.Client {
	return f.minioClient
}

func (f *ossItem) GetConfig() any {
	return f.config
}
