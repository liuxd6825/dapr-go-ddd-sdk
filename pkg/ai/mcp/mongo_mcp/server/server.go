package server

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/mongo_mcp/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/mongo_mcp/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Server struct {
	client    *mongo.Client
	mcpServer *server.MCPServer
	hServer   *server.StreamableHTTPServer
	cfg       *config.ServerConfig
	llm       llm.LLM
}

func (s *Server) AddTool(tool mcp.Tool, handler server.ToolHandlerFunc) {
	s.mcpServer.AddTool(tool, handler)
}

func NewServer(cfg *config.ServerConfig) *Server {
	ctx := context.Background()
	mcpServer := server.NewMCPServer("mongo-mcp", "1.0.0", server.WithToolCapabilities(false))
	client := connect(cfg)
	logger := logrus.New()
	openApiCfg := openai.ChatModelConfig{
		APIKey:  cfg.LlmApiKey,
		Model:   cfg.LlmModel,
		BaseURL: cfg.LlmBinding,
		Timeout: 500,
	}

	llm, err := llm.NewOpenAI(ctx, openApiCfg)
	if err != nil {
		panic(err)
	}

	s := &Server{client: client, cfg: cfg, mcpServer: mcpServer}
	for _, tool := range cfg.Tools {
		tools.NewMongoQueryTool(llm, client, &tool, logger).Register(s)
	}
	return s
}

func connect(cfg *config.ServerConfig) *mongo.Client {
	connectTimeout := time.Duration(cfg.ConnectTimeout)
	opt := &options.ClientOptions{
		AppName:        &cfg.AppName,
		ConnectTimeout: &connectTimeout,
		Auth: &options.Credential{
			Username:      cfg.UserName,
			Password:      cfg.Password,
			AuthSource:    cfg.AutoDB,
			AuthMechanism: cfg.AuthMechanism,
		},
	}
	ctx := context.Background()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return client
}

func (s *Server) Start() error {
	hServer := server.NewStreamableHTTPServer(s.mcpServer)
	s.hServer = hServer
	return hServer.Start(s.cfg.Addr)
}
