package tools

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/mcp/mongo_mcp/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

const SQLPrompt = `
Provided this schema:

CREATE TABLE mcp_item (
    id varchar(255) NOT NULL PRIMARY KEY, // 主键 
    case_id varchar(255) NOT NULL, // 案装ID
    tenant_id varchar(255) NOT NULL, // 租户ID
	name varchar(255) NOT NULL, // 姓名
    age int NOT NULL,       // 年龄
    money float NOT NULL,  // 金额
    date datetime NOT NULL, // 交易日期
);

Give me taxis with more than 2 passengers
`

type McpRecord struct {
	Id       string     `gorm:"id;primaryKey"`
	TenantId string     `gorm:"tenant_id"`
	Name     string     `gorm:"name"`
	Age      int        `gorm:"age"`
	Money    float64    `gorm:"money"`
	Date     *time.Time `gorm:"data"`
}

func Test_MongoQueryTool_Handler(t *testing.T) {
	envVal := xtest.NewEnvConfig_MongoLocalhost("db", "test")
	env.SetEnv(envVal)

	ctx := xtest.NewContext()
	llm := newOllamaLLM(ctx)

	logger := logrus.New()
	//inserts(ctx)

	db, ok := env.GetMongoByKey("db")
	if !ok {
		t.Fatal("db not exist")
	}
	cfg := &config.ToolConfig{
		Name:        "mcp_record",
		CollName:    "mcp_record",
		DBName:      "test",
		Description: "",
		SQLPrompt:   SQLPrompt,
	}
	tool := NewMongoQueryTool(llm, db.Client(), cfg, logger)
	params := map[string]any{
		"tenantId": "test",
		"caseId":   "1001",
		"userId":   "test",
		// "query":    "汇总年龄,平均金额，条件名称为张三， 日期在2020年",
		"query": "查询日期在2020年内的数据",
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "query",
			Params: mcp.RequestParams{
				Meta: &mcp.Meta{},
			},
		},
		Params: mcp.CallToolParams{
			Name:      "query",
			Arguments: params,
		},
	}

	resp, err := tool.handler(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range resp.Content {
		if val, ok := c.(mcp.TextContent); ok {
			t.Log(val.Text)
		}
	}

}

func inserts(ctx context.Context) {
	dao := dao.NewDao[*McpRecord](&dao.NewConfig{
		DBKey:     "db",
		TableName: "mcp_record",
	})

	var items []*McpRecord
	for i := 0; i < 10; i++ {
		item := &McpRecord{
			Name:  "张三",
			Date:  randomutils.PTime(),
			Money: 100,
		}
		items = append(items, item)
	}
	dao.CreateMany(ctx, items)
}

func newOllamaLLM(ctx context.Context) llm2.LLM {
	llm, err := llm2.NewOllama(ctx, ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "duckdb-nsql",
	})
	if err != nil {
		panic(err)
	}
	return llm
}
