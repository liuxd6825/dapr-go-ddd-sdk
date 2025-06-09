package common

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type IMCPServer interface {
	AddTool(tool mcp.Tool, handler server.ToolHandlerFunc)
}

type IMCPTool interface {
}
