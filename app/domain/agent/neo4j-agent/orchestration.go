package neo4jagent

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildGraphAgent(ctx context.Context) (r compose.Runnable[any, any], err error) {
	const (
		Text2CypherChatModel = "Text2CypherChatModel"
		Neo4jMcpToolsNode    = "Neo4jMcpToolsNode"
		DrawIoChatModel      = "DrawIoChatModel"
		ChatTemplate1        = "ChatTemplate1"
	)
	g := compose.NewGraph[any, any]()
	text2CypherChatModel, err := newText2CypherChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(Text2CypherChatModel, text2CypherChatModel)

	neo4jMcpToolsNode, err := newNeo4jMcpToolsNode(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddToolsNode(Neo4jMcpToolsNode, neo4jMcpToolsNode)

	drawIoChatModel, err := newOllamaChatModel(ctx)
	if err != nil {
		return nil, err
	}

	_ = g.AddChatModelNode(DrawIoChatModel, drawIoChatModel)
	chatTemplate1, err := newChatTemplate(ctx)

	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(ChatTemplate1, chatTemplate1)
	_ = g.AddEdge(compose.START, ChatTemplate1)
	_ = g.AddEdge(DrawIoChatModel, compose.END)
	_ = g.AddEdge(ChatTemplate1, Text2CypherChatModel)
	_ = g.AddEdge(Text2CypherChatModel, Neo4jMcpToolsNode)
	_ = g.AddEdge(Neo4jMcpToolsNode, DrawIoChatModel)
	_ = g.AddBranch(Neo4jMcpToolsNode, compose.NewGraphBranch(newBranch, map[string]bool{Text2CypherChatModel: true, Neo4jMcpToolsNode: true}))
	r, err = g.Compile(ctx, compose.WithGraphName("graphAgent"))
	if err != nil {
		return nil, err
	}
	return r, err
}
