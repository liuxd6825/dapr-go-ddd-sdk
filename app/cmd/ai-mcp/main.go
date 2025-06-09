package main

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/neo4j_mcp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/neo4j_mcp/server"
	"log"
)

func main() {
	cfg := &neo4j_mcp.Config{
		Name: "nc-neo4j-mcp",
		Addr: "0.0.0.0:6825",
		Neo4j: neo4j_mcp.Neo4jConfig{
			Uri:      "bolt://localhost:7687",
			Username: "neo4j",
			Password: "12345678",
		},
	}
	mcpServer := server.NewServer(cfg)
	err := mcpServer.Start()
	if err != nil {
		log.Fatalln(err)
	}
}
