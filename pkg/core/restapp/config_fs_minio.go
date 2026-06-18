package restapp

import (
	"context"
	"os"
	"strings"

	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

var _minioList = make(map[string]*minio.Client)
var _minioDefault *minio.Client

func InitMinioByEnvConfig(config *EnvConfig) error {
	if config != nil {
		return initMinio(config.Minio)
	}
	return nil
}

func initMinio(configs map[string]*MinioConfig) error {
	if err := assert2.NotNil(configs, assert2.NewOptions("config is nil")); err != nil {
		return err
	}
	ctx := context.Background()
	for k, c := range configs {
		if c.Endpoint == "<no value>" {
			logs.Errorf(ctx, nil, "MinIO  endpoint is not null in confg file")
			os.Exit(1)
		}
		k = strings.ToLower(k)
		options := &minio.Options{
			Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretKey, c.Token),
			Secure: c.UseSSL,
		}
		client, err := minio.New(c.Endpoint, options)
		if err != nil {
			logs.Errorf(ctx, nil, "连接MinIO:%s失败, error:%s.", c.Endpoint, err.Error())
			os.Exit(1)
		}
		c.client = client
		_minioList[k] = client
		if k == "default" {
			_minioDefault = client
		}
	}
	return nil
}

func GetMinioClient() *minio.Client {
	return _minioDefault
}

func SetMinioClient(miniClient *minio.Client) {
	_minioDefault = miniClient
}

func GetMinioClientByKey(key string) (*minio.Client, bool) {
	if len(key) == 0 {
		return _minioDefault, _minioDefault != nil
	}
	d, ok := _minioList[strings.ToLower(key)]
	return d, ok
}
