package restapp

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strings"
)

type MinioConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKey       string `yaml:"accessKey"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"useSSL"`
	minioClient     *minio.Client
}

var _minioList = make(map[string]*minio.Client)
var _minioDefault *minio.Client

func InitMinioByEnvConfig(config *EnvConfig) error {
	if config != nil {
		return InitMinio(config.Minio)
	}
	return nil
}

func InitMinio(configs map[string]*MinioConfig) error {
	if err := assert.NotNil(configs, assert.NewOptions("config is nil")); err != nil {
		return err
	}
	for k, c := range configs {
		if c.Endpoint == "<no value>" {
			continue
		}
		k = strings.ToLower(k)
		options := &minio.Options{
			Creds:  credentials.NewStaticV4(c.AccessKey, c.SecretAccessKey, ""),
			Secure: c.UseSSL,
		}
		minioClient, err := minio.New(c.Endpoint, options)
		if err != nil {
			return err
		}
		c.minioClient = minioClient
		_minioList[k] = minioClient
		if k == "default" {
			_minioDefault = minioClient
		}
	}
	return nil
}

func GetMinioClient() *minio.Client {
	return _minioDefault
}

func GetMinioClientByKey(key string) (*minio.Client, bool) {
	if len(key) == 0 {
		return _minioDefault, _mysqlDefault != nil
	}
	d, ok := _minioList[strings.ToLower(key)]
	return d, ok
}
