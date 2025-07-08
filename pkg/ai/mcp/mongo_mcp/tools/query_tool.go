package tools

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/mongo_mcp/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb/sql2mongo"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type MongoTool struct {
	tool       mcp.Tool
	cfg        *config.ToolConfig
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	llm        llm.LLM
	logger     *logrus.Logger
}

type QueryInParameter struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	UserId   string `json:"userId"`
	Query    string `json:"query"`
}

func NewMongoQueryTool(llm llm.LLM, client *mongo.Client, cfg *config.ToolConfig, logger *logrus.Logger) *MongoTool {
	db := client.Database(cfg.DBName)
	if db == nil {
		logger.Fatalf("failed to connect to database")
		panic("failed to connect to database")
	}
	collection := db.Collection(cfg.CollName)
	if collection == nil {
		logger.Fatalf("failed to connect to collection")
		panic("failed to connect to collection")
	}

	var opts []mcp.ToolOption
	if cfg.Description != "" {
		opts = append(opts, mcp.WithDescription(cfg.Description))
	}
	for _, item := range cfg.Options {
		var props []mcp.PropertyOption
		if item.Required {
			props = append(props, mcp.Required())
		}
		if item.Description != "" {
			props = append(props, mcp.Description(item.Description))
		}
		opt := mcp.WithString(item.Name, props...)
		opts = append(opts, opt)
	}

	tool := mcp.NewTool(cfg.Name, opts...)
	mongoTool := &MongoTool{
		client:     client,
		db:         db,
		llm:        llm,
		collection: collection,
		tool:       tool,
		cfg:        cfg,
		logger:     logger,
	}
	return mongoTool
}

func (q *MongoTool) getSQL(ctx context.Context, query string) (string, error) {
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(q.cfg.SQLPrompt),
		// 用户消息模板
		schema.UserMessage("问题: {query}"),
	)
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"query": query,
	})
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	resp, err := q.llm.Generate(ctx, messages)
	if err != nil {
		return "", err
	}
	sql := resp.Content
	if strings.HasSuffix(sql, ";") {
		sql = sql[:len(sql)-len(";")]
	}
	return sql, err
}

func (q *MongoTool) handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inParam, err := q.getRequest(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	q.logger.Infof("tenantId:%s; caseId:%s; userId:%s; query:%s", inParam.TenantId, inParam.CaseId, inParam.UserId, inParam.Query)
	sql, err := q.getSQL(ctx, inParam.Query)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	q.logger.Infof("sql:%s", sql)

	// 在此解析请求内容并执行 Neo4j 查询
	pipeline, err := sql2mongo.Pipeline(sql)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	cur, err := q.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var list []map[string]any
	err = cur.All(ctx, &list)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}

	if len(list) == 0 {
		return mcp.NewToolResultText("```json\n[]\n```"), nil
	}

	dataJson, err := json.Marshal(list)
	if err != nil {
		logrus.Error(err.Error())
		return mcp.NewToolResultError(err.Error()), nil
	}

	sb := new(strings.Builder)
	sb.WriteString("```json\n")
	sb.Write(dataJson)
	sb.WriteString("```")
	return mcp.NewToolResultText(sb.String()), nil
}

func (q *MongoTool) getRequest(request mcp.CallToolRequest) (*QueryInParameter, error) {
	tenantId, err := request.RequireString("tenantId")
	if err != nil {
		return nil, err
	}

	caseId, err := request.RequireString("caseId")
	if err != nil {
		return nil, err
	}

	userId, err := request.RequireString("userId")
	if err != nil {
		return nil, err
	}

	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}

	return &QueryInParameter{
		TenantId: tenantId,
		CaseId:   caseId,
		UserId:   userId,
		Query:    query,
	}, nil
}

func (q *MongoTool) Register(s common.IMCPServer) {
	s.AddTool(q.tool, q.handler)
}
