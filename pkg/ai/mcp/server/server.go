package server

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/tools/inter"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/tools/mongo"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/tools/neo4j"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

type Server struct {
	mcpServer *server.MCPServer
	hServer   *server.StreamableHTTPServer
	cfg       *Config
	logger    *logrus.Logger
}

func NewServer(cfg *Config, logger *logrus.Logger) *Server {
	tools := newTools(cfg, logger)
	mcpServer := server.NewMCPServer(cfg.Name, cfg.Version, server.WithToolCapabilities(false))
	s := &Server{cfg: cfg, mcpServer: mcpServer}
	for _, tool := range tools {
		s.mcpServer.AddTool(tool.GetTool(), tool.Handler)
	}
	return s
}

func newTools(cfg *Config, logger *logrus.Logger) []inter.Tool {
	var tools []inter.Tool
	ctx := context.Background()
	for _, item := range cfg.Mongo {
		tool, err := mongo.NewMongoQueryTool(ctx, &item, logger)
		if err != nil {
			panic(err)
		}
		tools = append(tools, tool)
	}
	for _, item := range cfg.Neo4j {
		tool, err := neo4j.NewNeo4jQueryTool(ctx, &item, logger)
		if err != nil {
			panic(err)
		}
		tools = append(tools, tool)
	}
	return tools
}

func (s *Server) Start() error {
	hServer := server.NewStreamableHTTPServer(s.mcpServer)
	s.hServer = hServer
	return hServer.Start(s.cfg.Addr)
}
