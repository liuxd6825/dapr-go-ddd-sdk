package restapp

import (
	"github.com/minio/minio-go/v7"
)

type MinioConfig struct {
	Name      string
	Endpoint  string            `json:"endpoint" yaml:"endpoint"`
	AccessKey string            `json:"accessKey" yaml:"accessKey"`
	SecretKey string            `json:"secretKey" yaml:"secretKey"`
	UseSSL    bool              `json:"useSSL" yaml:"useSSL"`
	Token     string            `json:"token" yaml:"token"`
	Buckets   map[string]string `json:"buckets" yaml:"buckets"` // 默认使用的存储桶名称
	client    *minio.Client
}
