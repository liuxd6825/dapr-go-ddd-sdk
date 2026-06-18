package env

import (
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	Name      string
	Endpoint  string            `yaml:"endpoint"`
	AccessKey string            `yaml:"accessKey"`
	SecretKey string            `yaml:"secretKey"`
	UseSSL    bool              `yaml:"useSSL"`
	Buckets   map[string]string `yaml:"buckets"`
	Client    *minio.Client
}

func NewMinio() *Minio {
	return &Minio{}
}

var _minioList = make(map[string]*minio.Client)

func initMinio(env *Env) {
	if env.Minio == nil {
		env.Minio = map[string]*Minio{}
		return
	}
	for k, m := range env.Minio {
		if m.Endpoint == "<no value>" {
			continue
		}
		k = strings.ToLower(k)
		options := &minio.Options{
			Creds:  credentials.NewStaticV4(m.AccessKey, m.SecretKey, ""),
			Secure: m.UseSSL,
		}
		client, err := minio.New(m.Endpoint, options)
		if err != nil {
			panic(fmt.Sprintf("minio连接%s失败, error:%s", m.Endpoint, err.Error()))
		}
		m.Client = client
		_minioList[k] = client
	}
}

func (m *Minio) Minio() {

}
