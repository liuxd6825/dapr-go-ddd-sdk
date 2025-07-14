package inter

import (
	"context"
	"github.com/mark3labs/mcp-go/mcp"
)

type Tool interface {
	GetTool() mcp.Tool
	Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}
