package envs

import (
	"fmt"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type Huge struct {
	DbKey    string `yaml:"dbKey" json:"dbKey"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Graph    string `yaml:"graph"`
	Username string `yaml:"user"`
	Password string `yaml:"pwd"`
	client   *hugegraph.CommonClient
}

func NewHuge() *Huge {
	return &Huge{
		Host:     "127.0.0.1",
		Port:     8080,
		Graph:    "hugegraph",
		Username: "admin",
		Password: "pa",
	}
}

func NewHugeWithMap(meta map[string]any) (*Huge, error) {
	var val Huge
	if err := copyObject(meta, &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func AddGraphDBHuge(env *Env, config *Huge) error {
	if config.Host == "<no value>" {
		return errors.New("db hugegraph host is empty")
	}
	if config.Graph == "" {
		config.Graph = "hugegraph"
	}
	client, err := hugegraph.NewCommonClient(hugegraph.Config{
		Host:     config.Host,
		Port:     config.Port,
		Graph:    config.Graph,
		Username: config.Username,
		Password: config.Password,
	})
	if err != nil {
		panic(fmt.Sprintf("hugegraph client error: %s", err.Error()))
	}
	res, err := client.Version()
	if err != nil {
		panic(fmt.Sprintf("connect hugegraph failed, error: %s", err.Error()))
	}
	_ = res.Body.Close()
	dbKey := strings.ToLower(config.DbKey)
	config.DbKey = dbKey
	env.AddDB(&dbItem{dbKey: dbKey, dbType: DBTypeHuge, hugeClient: client, config: config})
	return nil
}
