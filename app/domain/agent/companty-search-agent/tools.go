package companty_search_agent

import (
	"context"
	"github.com/cloudwego/eino-ext/components/tool/browseruse"
	"github.com/tmc/langchaingo/tools"
)

func NewTools(ctx context.Context) (tools.Tool, error) {
	tool, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	return tool, err
}
