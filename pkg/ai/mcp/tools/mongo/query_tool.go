package mongo

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/tools/inter"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type MongoQueryTool struct {
	tool       mcp.Tool
	cfg        *Config
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
	llm        llm.LLM
	logger     *logrus.Logger
	template   *prompt.DefaultChatTemplate
}

type QueryInParameter struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	UserId   string `json:"userId"`
	Query    string `json:"query"`
}

func NewMongoQueryTool(ctx context.Context, cfg *Config, logger *logrus.Logger) (inter.Tool, error) {
	return newMongoQueryTool(ctx, cfg, logger)
}

func newMongoQueryTool(ctx context.Context, cfg *Config, logger *logrus.Logger) (*MongoQueryTool, error) {
	llm, err := newLLM(ctx, &cfg.LLM)
	if err != nil {
		return nil, errors.ErrorOf("failed to connect to llm %s", err.Error())
	}

	client, err := getMongoDBClient(cfg)
	if err != nil {
		return nil, errors.ErrorOf("failed to connect to database %s", err.Error())
	}

	db := client.Database(cfg.DB.DBName)
	if db == nil {
		return nil, errors.ErrorOf("failed to connect to database %s", cfg.DB.DBName)
	}

	collection := db.Collection(cfg.DB.CollName)
	if collection == nil {
		return nil, errors.ErrorOf("failed to connect to collection %s ", cfg.DB.CollName)
	}

	var opts []mcp.ToolOption
	if cfg.Tool.Description != "" {
		opts = append(opts, mcp.WithDescription(cfg.Tool.Description))
	}
	if cfg.Tool.Prompt == "" {
		cfg.Tool.Prompt = PromptFormat
	}

	// 创建模板，使用 FString 格式
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(cfg.Tool.Prompt),
	)

	tool := newTool(cfg)

	mongoTool := &MongoQueryTool{
		client:     client,
		db:         db,
		llm:        llm,
		collection: collection,
		tool:       tool,
		cfg:        cfg,
		logger:     logger,
		template:   template,
	}
	return mongoTool, nil

}

func (q *MongoQueryTool) getPipeline(ctx context.Context, query string) (any, error) {
	// 使用模板生成消息
	messages, err := q.template.Format(ctx, map[string]any{
		"query":  query,
		"schema": q.cfg.Tool.Schema,
	})
	if err != nil {
		logrus.Fatal(err)
		return nil, errors.ErrorOf("使用模板生成消息提示词时出错：" + err.Error())
	}
	resp, err := q.llm.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	result := resp.Content

	i := strings.Index(result, "```json\n")
	if i < 0 {
		return nil, errors.New("没有找到json内容")
	}
	last := strings.LastIndex(result, "```")
	jsonData := result[i+7 : last]

	println(jsonData)
	var pip map[string]any
	err = json.Unmarshal([]byte(jsonData), &pip)
	if err != nil {
		return nil, err
	}
	if aggregate, ok := pip["aggregate"]; ok {
		return aggregate, nil
	}
	if pipeline, ok := pip["pipeline"]; ok {
		return pipeline, nil
	}
	return nil, errors.ErrorOf("没有找到aggregate,pipeline内容")
}

func (q *MongoQueryTool) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inParam, err := q.getRequest(request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	q.logger.Infof("tenantId:%s; caseId:%s; userId:%s; query:%s", inParam.TenantId, inParam.CaseId, inParam.UserId, inParam.Query)

	// 在此解析请求内容并执行 Neo4j 查询
	pipeline, err := q.getPipeline(ctx, inParam.Query)
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

func (q *MongoQueryTool) getRequest(request mcp.CallToolRequest) (*QueryInParameter, error) {
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

func newTool(cfg *Config) mcp.Tool {
	return mcp.NewTool(cfg.Tool.Name,
		mcp.WithDescription(cfg.Tool.Description),
		mcp.WithString("tenantId",
			mcp.Required(),
			mcp.Description("Tenant ID in the system"),
		),
		mcp.WithString("caseId",
			mcp.Required(),
			mcp.Description("ID of the project case"),
		),
		mcp.WithString("userId",
			mcp.Required(),
			mcp.Description("User ID in the system"),
		),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The natural language to be queried"),
		))
}
