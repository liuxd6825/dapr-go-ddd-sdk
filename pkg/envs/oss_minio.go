package envs

import (
	"fmt"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	OssKey    string
	Endpoint  string            `yaml:"endpoint"`
	AccessKey string            `yaml:"accessKey"`
	SecretKey string            `yaml:"secretKey"`
	UseSSL    bool              `yaml:"useSSL"`
	Buckets   map[string]string `yaml:"buckets"`
	Client    *minio.Client
}

func newMinioWithMap(meta Metadata) (*Minio, error) {
	var val Minio
	if err := copyObject(meta, &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func addOssMinio(env *Env, config *Minio) error {
	k := strings.ToLower(config.OssKey)
	if config.Endpoint == "<no value>" {
		return errors.New("no endpoint")
	}
	options := &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	}
	client, err := minio.New(config.Endpoint, options)
	if err != nil {
		panic(fmt.Sprintf("minio连接%s失败, error:%s", config.Endpoint, err.Error()))
	}
	config.Client = client
	env.AddOSS(&ossItem{fsKey: k, ossType: OssTypeMinio, minioClient: client, config: config})
	return nil
}
