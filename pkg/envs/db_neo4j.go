package envs

import (
	"context"
	"fmt"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4j struct {
	DbKey        string `yaml:"dbKey" json:"dbKey"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"dbname"`
	UserName     string `yaml:"user"`
	Password     string `yaml:"password"`
	EventPublish bool   `yaml:"eventPublish" json:"eventPublish"` // 是否发送领域事件
}

func NewNeo4j() *Neo4j {
	return &Neo4j{}
}

func NewNeo4jWithMap(meta map[string]any) (*Neo4j, error) {
	var val Neo4j
	if err := copyObject(meta, &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func AddDBNeo4j(env *Env, config *Neo4j) error {
	if config.Host == "<no value>" && config.Port == "<no value>" {
		return errors.New("")
	}
	uri := fmt.Sprintf("bolt://%v:%v", config.Host, config.Port)
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(config.UserName, config.Password, ""))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(fmt.Sprintf("连接neo4j失败, error:%s。%s  ", uri, err.Error()))
	}
	dbKey := strings.ToLower(config.DbKey)
	config.DbKey = dbKey
	env.AddDB(&dbItem{dbKey: dbKey, dbType: DBTypeNeo4j, neo4j: driver, config: config})
	return nil
}
