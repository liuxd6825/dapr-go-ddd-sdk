package server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/tools/mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/tools/neo4j"
)

type Config struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Addr    string `json:"addr" yaml:"addr"`
	// Tools   []ToolConfig `json:"tools" yaml:"tools"`
	Mongo []mongo.Config `json:"mongo" yaml:"mongo"`
	Neo4j []neo4j.Config `json:"neo4j" yaml:"neo4J"`
}

type ToolConfig struct {
	Kind string `json:"kind" yaml:"kind"`
	File string `json:"file" yaml:"file"`
}
