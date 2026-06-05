package env

import (
	"context"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

type Elastic struct {
	DbKey    string `json:"dbKey"`
	Name     string `yaml:"name" json:"name"`
	Addresses    []string `yaml:"addresses" json:"addresses"`
	Username     string  `yaml:"username" json:"username"`
	Password     string  `yaml:"password" json:"password"`
	APIKey       string  `yaml:"apiKey" json:"apiKey"`
	ServiceToken string  `yaml:"serviceToken" json:"serviceToken"`
	CloudID      string  `yaml:"cloudId" json:"cloudId"`
	CACert       string  `yaml:"caCert" json:"caCert"`
	ClientCert   string  `yaml:"clientCert" json:"clientCert"`
	ClientKey    string  `yaml:"clientKey" json:"clientKey"`
	MaxRetries   int     `yaml:"maxRetries" json:"maxRetries"`
	Debug        bool    `yaml:"debug" json:"debug"`
}

func NewElastic() *Elastic {
	return &Elastic{}
}

func (e *Elastic) IsEmpty() bool {
	if len(e.Addresses) == 0 && e.CloudID == "" && e.Username == "" && e.Password == "" && e.APIKey == "" && e.ServiceToken == "" {
		return true
	}
	return false
}

func InitDBElastic(env *Env) {
	if env.Elastic == nil {
		env.Elastic = map[string]*Elastic{}
		return
	}

	for k, ec := range env.Elastic {
		if ec.IsEmpty() {
			continue
		}

		var opts []elasticsearch.Option
		opts = append(opts, elasticsearch.WithAddresses(ec.Addresses...))

		if ec.CloudID != "" {
			opts = append(opts, elasticsearch.WithCloudID(ec.CloudID))
		}

		switch {
		case ec.APIKey != "":
			opts = append(opts, elasticsearch.WithAPIKey(ec.APIKey))
		case ec.ServiceToken != "":
			opts = append(opts, elasticsearch.WithServiceToken(ec.ServiceToken))
		case ec.Username != "" && ec.Password != "":
			opts = append(opts, elasticsearch.WithBasicAuth(ec.Username, ec.Password))
		}

		if ec.CACert != "" {
			opts = append(opts, elasticsearch.WithCACert([]byte(ec.CACert)))
		}

		if ec.MaxRetries > 0 {
			statuses := []int{502, 503, 504}
			opts = append(opts, elasticsearch.WithRetry(ec.MaxRetries, statuses...))
		}

		client, err := elasticsearch.New(opts...)
		if err != nil {
			panic(fmt.Sprintf("elastic连接%s失败, error:%s", ec.Addresses[0], err.Error()))
		}

		logs.Infofmt(context.Background(), "config elastic addresses=%v; username=%s; cloudId=%s",
			ec.Addresses, ec.Username, ec.CloudID)

		dbKey := strings.ToLower(k)
		ec.DbKey = dbKey
		env.AddDB(&dbItem{
			dbKey:   dbKey,
			dbType:  DBType_Elastic,
			elastic: client,
			config:  ec,
		})
	}
}