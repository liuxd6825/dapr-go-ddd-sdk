package server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/neo4j_mcp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/neo4j_mcp/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type Server struct {
	driver    neo4j.Driver
	mcpServer *server.MCPServer
	hServer   *server.StreamableHTTPServer
	cfg       *neo4j_mcp.Config
}

func (s *Server) AddTool(tool mcp.Tool, handler server.ToolHandlerFunc) {
	s.mcpServer.AddTool(tool, handler)
}

func NewServer(cfg *neo4j_mcp.Config) *Server {
	mcpServer := server.NewMCPServer("neo4j-mcp", "1.0.0", server.WithToolCapabilities(false))
	driver := connect(cfg)
	s := &Server{driver: driver, cfg: cfg, mcpServer: mcpServer}
	tools.NewNeo4jQueryTool(driver).Register(s)
	return s
}

func connect(cfg *neo4j_mcp.Config) neo4j.Driver {
	driver, err := neo4j.NewDriver(cfg.Neo4j.Uri, neo4j.BasicAuth(cfg.Neo4j.Username, cfg.Neo4j.Password, ""))
	if err != nil {
		log.Fatalf("failed to create the driver: %v", err)
	}
	return driver
}

func (s *Server) Start() error {
	hServer := server.NewStreamableHTTPServer(s.mcpServer)
	s.hServer = hServer
	return hServer.Start(s.cfg.Addr)
}
