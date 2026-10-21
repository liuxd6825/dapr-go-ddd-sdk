package xtest

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"

const (
	Neo4jDBKey      = "db"
	Neo4jDBName     = ""
	Neo4jHostLocal  = "localhost"
	Neo4jHostRemote = "192.168.120.224"
)

func GetNeo4jRemoteCfg() *envs.Neo4j {
	return &envs.Neo4j{
		DbKey:    Neo4jDBKey,
		Host:     Neo4jHostRemote,
		Port:     "7687",
		Database: Neo4jDBName,
		UserName: "neo4jdao",
		Password: "12345678",
	}
}

func GetNeo4jLocalCfg() *envs.Neo4j {
	return &envs.Neo4j{
		DbKey:    Neo4jDBKey,
		Host:     Neo4jHostLocal,
		Port:     "7687",
		Database: "",
		UserName: "neo4jdao",
		Password: "12345678",
	}
}
