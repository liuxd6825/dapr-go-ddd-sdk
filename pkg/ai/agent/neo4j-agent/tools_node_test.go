package neo4jagent

import (
	"context"
	tool_mcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/agent/common"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"testing"
)

func Test_newTools(t *testing.T) {
	ctx := context.Background()
	// Mock LLM 输出作为输入
	input := &schema.Message{
		Role: schema.Assistant,
		ToolCalls: []schema.ToolCall{
			{
				Function: schema.FunctionCall{
					Name:      "weather",
					Arguments: `{"city": "深圳", "date": "tomorrow"}`,
				},
			},
		},
	}

	chatModel, err := newOllamaChatModel(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = chatModel.WithTools(nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	reader, err := chatModel.Stream(ctx, []*schema.Message{input})
	if err != nil {
		t.Fatal(err)
		return
	}

	sb, err := common.Reader(reader)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(sb.String())
}

func Test_callMcpTool(t *testing.T) {
	ctx := context.Background()

	cm, err := newText2CypherChatModel(ctx)
	if err != nil {
		t.Fatal(err)
	}

	mcpTools, err := newMcpTools(ctx, t)
	if err != nil {
		t.Fatal(err)
	}
	conf := &compose.ToolsNodeConfig{
		Tools: mcpTools,
	}
	// 创建工具节点
	toolsNode, err := compose.NewToolNode(ctx, conf)

	chatTpl := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage(`
你是一个关系数据查询器,将自然语言转换为neo4j的Cypher语句,将调用MCP查询关系数据，
参数tenantId为{tenantId}, 参数caseId为{caseId}, 参数userId为{userId}。
人员的标签是：human, 名称属性：name;
账号的标签是：account，名称属性：name;
生成类似cypher语句：MATCH p=(a:account)-[*]-(e:human ) RETURN p ;
`),
		schema.UserMessage("{query}"),
	)

	mcpToolInfo := []*schema.ToolInfo{
		{
			Name: "neo4j-query",
			Desc: "通过neo4j的cypher语句查询关系数据",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"tenantId": {
					Type:     "string",
					Required: true,
					Desc:     "租户id",
				},
				"caseId": {
					Type:     "string",
					Required: true,
					Desc:     "案件id",
				},
				"userId": {
					Type:     "string",
					Required: true,
					Desc:     "用户id",
				},
				"cypher": {
					Type:     "string",
					Required: true,
					Desc:     "neo4j的cypher语句",
				},
			}),
		},
	}
	cmWithTools, err := cm.WithTools(mcpToolInfo)
	if err != nil {
		log.Fatalf("Failed to bind tools to ChatModel: %v", err)
	}

	// 在 Chain 中使用
	chain := compose.NewChain[map[string]interface{}, []*schema.Message]()
	chain.AppendChatTemplate(chatTpl)
	chain.AppendChatModel(cmWithTools)
	chain.AppendToolsNode(toolsNode)

	c, err := chain.Compile(ctx)
	if err != nil {
		t.Fatal(err)
	}

	input := map[string]interface{}{
		"tenantId": "test",
		"caseId":   "1001",
		"userId":   "user",
		"query":    "查找账号0001与人员王二的全关系路径",
	}
	output, err := c.Invoke(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range output {
		t.Log(o.Content)
	}

}

func newMcpTools(ctx context.Context, t *testing.T) ([]tool.BaseTool, error) {
	cli, err := client.NewStreamableHttpClient("http://127.0.0.1:6825/mcp")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "eino-client",
		Version: "1.0.0",
	}
	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatal(err)
	}

	mcpTools, err := tool_mcp.GetTools(ctx, &tool_mcp.Config{Cli: cli})
	return mcpTools, err
}
