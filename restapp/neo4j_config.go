package restapp

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"strings"
)

type Neo4jConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"dbname"`
	UserName string `yaml:"user"`
	Password string `yaml:"pwd"`
}

var _neo4js = make(map[string]neo4j.Driver)
var _neo4jDefault neo4j.Driver

func InitNeo4j(configs map[string]*Neo4jConfig) {
	if err := assert.NotNil(configs, assert.NewOptions("config is nil")); err != nil {
		panic(err)
	}

	for key, config := range configs {
		if config.Host == "<no value>" && config.Port == "<no value>" {
			continue
		}
		uri := fmt.Sprintf("bolt://%v:%v", config.Host, config.Port)
		driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(config.UserName, config.Password, ""))
		if err != nil {
			panic(err)
		}
		_neo4js[strings.ToLower(key)] = driver
		_neo4jDefault = driver
	}
}

func GetNeo4j() neo4j.Driver {
	return _neo4jDefault
}

func GetNeo4jByKey(dbKey string) (neo4j.Driver, bool) {
	if len(dbKey) == 0 {
		return _neo4jDefault, _neo4jDefault != nil
	}
	d, ok := _neo4js[strings.ToLower(dbKey)]
	return d, ok
}

func CloseAllNeo4j(ctx context.Context) error {
	c := func(d neo4j.Driver) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		return d.Close()
	}
	for _, d := range _neo4js {
		_ = c(d)
	}
	return nil
}
