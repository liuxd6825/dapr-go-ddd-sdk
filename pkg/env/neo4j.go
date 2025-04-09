package env

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strings"
)

type Neo4j struct {
	DbKey        string
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"dbname"`
	UserName     string `yaml:"user"`
	Password     string `yaml:"pwd"`
	EventPublish bool   `yaml:"eventPublish" json:"eventPublish"` // 是否发送领域事件
}

func NewNeo4j() *Neo4j {
	return &Neo4j{}
}

var _neo4js = make(map[string]neo4j.DriverWithContext)
var _neo4jDefault neo4j.DriverWithContext

func initNeo4j(env *Env) {
	if env.Neo4j == nil {
		env.Neo4j = map[string]*Neo4j{}
		return
	}

	ctx := context.Background()
	for dbKey, config := range env.Neo4j {
		if config.Host == "<no value>" && config.Port == "<no value>" {
			continue
		}
		uri := fmt.Sprintf("bolt://%v:%v", config.Host, config.Port)
		driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(config.UserName, config.Password, ""))
		if err != nil {
			panic(err)
		}
		err = driver.VerifyConnectivity(ctx)
		if err != nil {
			panic(fmt.Sprintf("连接neo4j失败, error:%s。%s  ", uri, err.Error()))
		}
		key := strings.ToLower(dbKey)
		config.DbKey = key
		env.AddDB(&dbItem{dbKey: dbKey, dbType: DBType_Neo4j, neo4j: driver, config: config})
	}
}
