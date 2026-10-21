package envs

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
)

type Huge struct {
	DbKey    string
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

var _huges = make(map[string]*hugegraph.CommonClient)

func InitDBHuge(env *Env) {
	if env.Huge == nil {
		env.Huge = map[string]*Huge{}
		return
	}

	ctx := context.Background()
	_ = ctx
	for dbKey, config := range env.Huge {
		if config.Host == "<no value>" {
			continue
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
		key := strings.ToLower(dbKey)
		config.DbKey = key
		env.AddDB(&dbItem{dbKey: dbKey, dbType: DBType_Huge, hugeClient: client, config: config})
	}
}

func (env *Env) GetHugeByKey(dbKey string) (*Huge, bool) {
	if env.Huge == nil {
		return nil, false
	}
	n := env.Huge[dbKey]
	return n, n != nil
}
